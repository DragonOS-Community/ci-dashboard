package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/dragonos/dragonos-ci-dashboard/internal/services"
	"golang.org/x/term"
)

var (
	action      = flag.String("action", "", "操作类型: create, update-password, update-role")
	username    = flag.String("username", "", "用户名")
	password    = flag.String("password", "", "密码（如果不提供，将提示输入）")
	role        = flag.String("role", "admin", "角色: admin 或 user")
	interactive = flag.Bool("interactive", false, "交互式模式")
)

func main() {
	flag.Parse()

	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := models.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDatabase()

	// 交互式模式
	if *interactive {
		runInteractive()
		return
	}

	// 命令行模式
	if *action == "" {
		fmt.Println("Usage: admin-cli [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  admin-cli -action=create -username=admin -password=secret123")
		fmt.Println("  admin-cli -action=update-password -username=admin")
		fmt.Println("  admin-cli -action=update-role -username=admin -role=user")
		fmt.Println("  admin-cli -interactive")
		os.Exit(1)
	}

	switch *action {
	case "create":
		handleCreate()
	case "update-password":
		handleUpdatePassword()
	case "update-role":
		handleUpdateRole()
	default:
		log.Fatalf("Unknown action: %s. Use: create, update-password, update-role", *action)
	}
}

func runInteractive() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== DragonOS CI Dashboard Admin CLI ===")
		fmt.Println("1. 创建管理员账户")
		fmt.Println("2. 更改管理员密码")
		fmt.Println("3. 更改管理员角色")
		fmt.Println("4. 退出")
		fmt.Print("\n请选择操作 (1-4): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			interactiveCreate(reader)
		case "2":
			interactiveUpdatePassword(reader)
		case "3":
			interactiveUpdateRole(reader)
		case "4":
			fmt.Println("退出")
			return
		default:
			fmt.Println("无效的选择，请重试")
		}
	}
}

func interactiveCreate(reader *bufio.Reader) {
	fmt.Print("用户名: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		fmt.Println("用户名不能为空")
		return
	}

	password := readPassword("密码: ")
	if password == "" {
		fmt.Println("密码不能为空")
		return
	}

	fmt.Print("角色 (admin/user, 默认: admin): ")
	roleInput, _ := reader.ReadString('\n')
	roleInput = strings.TrimSpace(roleInput)
	if roleInput == "" {
		roleInput = "admin"
	}

	role := models.UserRole(roleInput)
	if role != models.UserRoleAdmin && role != models.UserRoleUser {
		fmt.Println("无效的角色，使用默认值: admin")
		role = models.UserRoleAdmin
	}

	user, err := services.CreateUser(username, password, role)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 用户创建成功!\n")
	fmt.Printf("  ID: %d\n", user.ID)
	fmt.Printf("  用户名: %s\n", user.Username)
	fmt.Printf("  角色: %s\n", user.Role)
}

func interactiveUpdatePassword(reader *bufio.Reader) {
	fmt.Print("用户名: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		fmt.Println("用户名不能为空")
		return
	}

	password := readPassword("新密码: ")
	if password == "" {
		fmt.Println("密码不能为空")
		return
	}

	err := services.UpdateUserPassword(username, password)
	if err != nil {
		fmt.Printf("更新密码失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 密码更新成功!\n")
}

func interactiveUpdateRole(reader *bufio.Reader) {
	fmt.Print("用户名: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		fmt.Println("用户名不能为空")
		return
	}

	fmt.Print("新角色 (admin/user): ")
	roleInput, _ := reader.ReadString('\n')
	roleInput = strings.TrimSpace(roleInput)

	role := models.UserRole(roleInput)
	if role != models.UserRoleAdmin && role != models.UserRoleUser {
		fmt.Println("无效的角色，必须是 admin 或 user")
		return
	}

	err := services.UpdateUserRole(username, role)
	if err != nil {
		fmt.Printf("更新角色失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 角色更新成功!\n")
}

func handleCreate() {
	if *username == "" {
		log.Fatal("Username is required for create action")
	}

	password := *password
	if password == "" {
		password = readPassword("Password: ")
		if password == "" {
			log.Fatal("Password is required")
		}
	}

	role := models.UserRole(*role)
	if role != models.UserRoleAdmin && role != models.UserRoleUser {
		log.Fatal("Role must be 'admin' or 'user'")
	}

	user, err := services.CreateUser(*username, password, role)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("User created successfully!\n")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Role: %s\n", user.Role)
}

func handleUpdatePassword() {
	if *username == "" {
		log.Fatal("Username is required for update-password action")
	}

	password := *password
	if password == "" {
		password = readPassword("New password: ")
		if password == "" {
			log.Fatal("Password is required")
		}
	}

	err := services.UpdateUserPassword(*username, password)
	if err != nil {
		log.Fatalf("Failed to update password: %v", err)
	}

	fmt.Printf("Password updated successfully for user: %s\n", *username)
}

func handleUpdateRole() {
	if *username == "" {
		log.Fatal("Username is required for update-role action")
	}

	role := models.UserRole(*role)
	if role != models.UserRoleAdmin && role != models.UserRoleUser {
		log.Fatal("Role must be 'admin' or 'user'")
	}

	err := services.UpdateUserRole(*username, role)
	if err != nil {
		log.Fatalf("Failed to update role: %v", err)
	}

	fmt.Printf("Role updated successfully for user: %s\n", *username)
	fmt.Printf("New role: %s\n", role)
}

// readPassword 安全地读取密码（不显示在终端）
func readPassword(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println() // 换行
	return strings.TrimSpace(string(bytePassword))
}
