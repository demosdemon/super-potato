package main

var enums = []Enum{
	{
		Name: "AccessLevel",
		Values: []EnumValue{
			{
				Name:  "Viewer",
				Value: "viewer",
			},
			{
				Name:  "Contributor",
				Value: "contributor",
			},
			{
				Name:  "Admin",
				Value: "admin",
			},
		},
	},
	{
		Name: "AccessType",
		Values: []EnumValue{
			{
				Name:  "SSH",
				Value: "ssh",
			},
		},
	},
	{
		Name: "ApplicationMount",
		Values: []EnumValue{
			{
				Name:  "Local",
				Value: "local",
			},
			{
				Name:  "Temp",
				Value: "tmp",
			},
			{
				Name:  "Service",
				Value: "service",
			},
		},
	},
	{
		Name: "ServiceSize",
		Values: []EnumValue{
			{
				Name:  "Auto",
				Value: "AUTO",
			},
			{
				Name:  "Small",
				Value: "S",
			},
			{
				Name:  "Medium",
				Value: "M",
			},
			{
				Name:  "Large",
				Value: "L",
			},
			{
				Name:  "ExtraLarge",
				Value: "XL",
			},
			{
				Name:  "DoubleExtraLarge",
				Value: "2XL",
			},
			{
				Name:  "QuadrupleExtraLarge",
				Value: "4XL",
			},
		},
	},
	{
		Name: "SocketFamily",
		Values: []EnumValue{
			{
				Name:  "TCP",
				Value: "tcp",
			},
			{
				Name:  "Unix",
				Value: "unix",
			},
		},
	},
	{
		Name: "SocketProtocol",
		Values: []EnumValue{
			{
				Name:  "HTTP",
				Value: "http",
			},
			{
				Name:  "FastCGI",
				Value: "fastcgi",
			},
			{
				Name:  "UWSGI",
				Value: "uwsgi",
			},
		},
	},
}
