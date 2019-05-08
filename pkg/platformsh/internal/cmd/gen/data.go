package main

import (
	. "github.com/demosdemon/super-potato/pkg/platformsh/internal/cmd/gen/enums"
	. "github.com/demosdemon/super-potato/pkg/platformsh/internal/cmd/gen/env"
)

var enumData = Collection{
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

var wellKnownVariables = WellKnownVariables{
	{
		Name:           "Application",
		DecodedType:    "Application",
		DecodedPointer: true,
	},
	{
		Name:    "ApplicationName",
		Aliases: []string{"AppName"},
	},
	{
		Name:    "AppCommand",
		Aliases: []string{"ApplicationCommand"},
	},
	{
		Name: "AppDir",
	},
	{
		Name: "Branch",
	},
	{
		Name: "Dir",
	},
	{
		Name: "DocumentRoot",
	},
	{
		Name: "Environment",
	},
	{
		Name:     "Port",
		NoPrefix: true,
	},
	{
		Name: "Project",
	},
	{
		Name: "ProjectEntropy",
	},
	{
		Name:        "Relationships",
		DecodedType: "Relationships",
	},
	{
		Name:        "Routes",
		DecodedType: "Routes",
	},
	{
		Name: "SMTPHost",
	},
	{
		Name:     "Socket",
		NoPrefix: true,
	},
	{
		Name: "TreeID",
	},
	{
		Name:        "Variables",
		Aliases:     []string{"Vars"},
		DecodedType: "JSONObject",
	},
}
