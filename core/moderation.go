// Copyright (c) 2017-2019 Andrew Goulas
// Licensed under the MIT license.

package main

import (
	"net"
	"strings"

	"github.com/structinf/Go-MCC/gomcc"
)

func (plugin *Plugin) handleBan(sender gomcc.CommandSender, command *gomcc.Command, message string) {
	if len(message) == 0 {
		sender.SendMessage("Usage: " + command.Name + " <player> <reason>")
		return
	}

	reason := "You have been banned"
	args := strings.SplitN(message, " ", 2)
	if len(args) > 1 {
		reason = args[1]
	}

	if !gomcc.IsValidName(args[0]) {
		sender.SendMessage(args[0] + " is not a valid name")
		return
	}

	plugin.db.ban(args[0], reason, sender.Name())
	if player := sender.Server().FindPlayer(args[0]); player != nil {
		player.Kick(reason)
	}

	sender.SendMessage("Player " + args[0] + " banned")
}

func (plugin *Plugin) handleBanIp(sender gomcc.CommandSender, command *gomcc.Command, message string) {
	if len(message) == 0 {
		sender.SendMessage("Usage: " + command.Name + " <ip> <reason>")
		return
	}

	reason := "You have been banned"
	args := strings.SplitN(message, " ", 2)
	if len(args) > 1 {
		reason = args[1]
	}

	if net.ParseIP(args[0]) == nil {
		sender.SendMessage(args[0] + " is not a valid IP address")
		return
	}

	plugin.db.banIp(args[0], reason, sender.Name())
	sender.Server().ForEachPlayer(func(player *gomcc.Player) {
		if player.RemoteAddr() == args[0] {
			player.Kick(reason)
		}
	})

	sender.SendMessage("IP " + args[0] + " banned")
}

func (plugin *Plugin) handleKick(sender gomcc.CommandSender, command *gomcc.Command, message string) {
	if len(message) == 0 {
		sender.SendMessage("Usage: " + command.Name + " <player> <reason>")
		return
	}

	args := strings.SplitN(message, " ", 2)
	player := sender.Server().FindPlayer(args[0])
	if player == nil {
		sender.SendMessage("Player " + args[0] + " not found")
		return
	}

	reason := "Kicked by " + sender.Name()
	if len(args) > 1 {
		reason = args[1]
	}

	player.Kick(reason)
}

func (plugin *Plugin) handleRank(sender gomcc.CommandSender, command *gomcc.Command, message string) {
	var rank *gomcc.Rank
	args := strings.Fields(message)
	switch len(args) {
	case 1:
		rank = nil

	case 2:
		if rank = plugin.findRank(args[1]); rank == nil {
			sender.SendMessage("Rank " + args[1] + " not found")
			return
		}

	default:
		sender.SendMessage("Usage: " + command.Name + " <player> <rank>")
		return
	}

	if player := plugin.findPlayer(args[0]); player == nil {
		sender.SendMessage("Player " + args[0] + " not found")
	} else {
		player.Rank = rank
		player.SendPermissions()
		if rank == nil {
			sender.SendMessage("Rank of " + args[0] + " reset")
		} else {
			sender.SendMessage("Rank of " + args[0] + " set to " + args[1])
		}
	}
}

func (plugin *Plugin) handleUnban(sender gomcc.CommandSender, command *gomcc.Command, message string) {
	args := strings.Fields(message)
	if len(args) != 1 {
		sender.SendMessage("Usage: " + command.Name + " <player>")
		return
	}

	if plugin.db.unban(args[0]) {
		sender.SendMessage("Player " + args[0] + " unbanned")
	} else {
		sender.SendMessage("Player " + args[0] + " is not banned")
	}
}

func (plugin *Plugin) handleUnbanIp(sender gomcc.CommandSender, command *gomcc.Command, message string) {
	args := strings.Fields(message)
	if len(args) != 1 {
		sender.SendMessage("Usage: " + command.Name + " <ip>")
		return
	}

	if plugin.db.unbanIp(args[0]) {
		sender.SendMessage("IP " + args[0] + " unbanned")
	} else {
		sender.SendMessage("IP " + args[0] + " is not banned")
	}
}
