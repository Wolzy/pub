数据库连接
支持Oracle/PostgreSQL/MySQL

配置文件示例：

	"MainDB": {
		"Name": "核心库",
		"Protocol": "MySQL",   # 目前可选MySQL/Oracle/PostgreSQL/PGSQL
		"IP": "10.242.31.180", # 必填
		"Port": "3306",        # 必填
		"DBName": "hibss",     # 必填
		"Network": "tcp",      # tcp/udp...
		"UserName": "hibss",   # 
		"PassWord": "Hibss135!",
		"NLS_LANG": "SIMPLIFIED CHINESE_CHINA.ZHS16GBK"
	}
