import { getLogSettings } from "../lib/getLogSettings";
import sendMessage from "../lib/sendMessage";
import logger from "../lib/logger";
import { time } from "discord.js";

// get memeber event log settings
async function getMemberEventLogSettings(logType: string, serverid: any) {
    const logSettings = await getLogSettings(logType, serverid);

    if (logSettings) {
        return logSettings;
    } else {
        return false;
    }
}

// member ban event
async function memberBanEvent(data: any) {
    const eventType = "member_banned";

    const event = data.event,
        guild = data.guild,
        user = data.user,
        reason = data.reason,
        moderator = data.moderator;

    let settings: any = await getMemberEventLogSettings(eventType, guild.id);
    if (!settings) {
        return;
    }

    settings = settings.settings;

    if (!settings.types.members) {
        return;
    }

    const embed = {
        title: "Member Banned",
        description: `**Member:** <@${user.id}> (${user.id})\n**Reason:** \`${reason}\`\n**Moderator:** <@${moderator}> (${moderator})`,
        thumbnail: {
            url: "https://cdn.discordapp.com/emojis/1064442673806704672.webp"
        },
        color: 0xe75151,
        author: {
            name: user.username,
            icon_url: user.avatarURL
        },
        footer: {
            text: "Event ID: " + event.id + " | " + eventType + " event"
        },
        timestamp: new Date()
    };

    const send = await sendMessage(
        {
            embeds: [embed]
        },
        settings.types.members.webhook_url
    ).catch((error) => {
        logger.error("Error sending " + eventType + " webhook", error);
        return { error: "Error sending " + eventType + " webhook" };
    });

    if (!send) {
        return { error: "Error sending " + eventType + " webhook" };
    }

    return true;
}

// member unban event
async function memberUnbanEvent(data: any) {
    const eventType = "member_unbanned";

    const event = data.event,
        guild = data.guild,
        user = data.user,
        reason = data.reason,
        moderator = data.moderator;

    let settings: any = await getMemberEventLogSettings(eventType, guild.id);
    if (!settings) {
        return;
    }

    settings = settings.settings;

    if (!settings.types.members) {
        return;
    }

    const embed = {
        title: "Member Unbanned",
        description: `**Member:** <@${user.id}> (${user.id})\n**Reason:** \`${reason}\`\n**Moderator:** <@${moderator}> (${moderator})`,
        thumbnail: {
            url: "https://cdn.discordapp.com/emojis/1064442704936828968.webp"
        },
        color: 0x2c2f33,
        author: {
            name: user.username,
            icon_url: user.avatarURL
        },
        footer: {
            text: "Event ID: " + event.id + " | " + eventType + " event"
        },
        timestamp: new Date()
    };

    const send = await sendMessage(
        {
            embeds: [embed]
        },
        settings.types.members.webhook_url
    ).catch((error) => {
        logger.error("Error sending " + eventType + " webhook", error);
        return { error: "Error sending " + eventType + " webhook" };
    });

    if (!send) {
        return { error: "Error sending " + eventType + " webhook" };
    }

    return true;
}

// member join event
async function memberJoinEvent(data: any) {
    const eventType = "member_join";

    const event = data.event,
        guild = data.guild,
        user = data.user;

    let settings: any = await getMemberEventLogSettings(eventType, guild.id);
    if (!settings) {
        return;
    }

    settings = settings.settings;

    if (!settings.types.members) {
        return;
    }

    const newAccount = Date.now() - user.createdTimestamp < 1000 * 60 * 60 * 24 * 7;

    const embed = {
        title: "Member Joined",
        description: `**Member:** <@${user.id}> (${user.id}) - Member #${guild.memberCount.toLocaleString(
            "en-US"
        )}\n**Account Created:** ${time(user.createdAt, "F")} ${newAccount ? "⚠️ New Account ⚠️" : ""}`,
        thumbnail: {
            url: "https://cdn.discordapp.com/emojis/1064442704936828968.webp"
        },
        color: 0xeb459e,
        author: {
            name: user.username,
            icon_url: user.avatarURL
        },
        footer: {
            text: "Event ID: " + event.id + " | " + eventType + " event"
        },
        timestamp: new Date()
    };

    const send = await sendMessage(
        {
            embeds: [embed]
        },
        settings.types.members.webhook_url
    ).catch((error) => {
        logger.error("Error sending " + eventType + " webhook", error);
        return { error: "Error sending " + eventType + " webhook" };
    });

    if (!send) {
        return { error: "Error sending " + eventType + " webhook" };
    }

    return true;
}

// member leave event
async function memberLeaveEvent(data: any) {
    const eventType = "member_leave";

    const event = data.event,
        guild = data.guild,
        user = data.user,
        roles = data.roles;

    let settings: any = await getMemberEventLogSettings(eventType, guild.id);
    if (!settings) {
        return;
    }

    settings = settings.settings;

    if (!settings.types.members) {
        return;
    }

    // get the users roles
    const rolesString = roles
        .map((role: any) => {
            return `<@&${role}>`;
        })
        .join(", ")
        .replace(`<@&${guild.id}>`, "");

    const embed = {
        title: "Member Left",
        description:
            "**Member:** <@" +
            user.id +
            "> (" +
            user.id +
            ") \n**Joined:** " +
            time(user.joinedAt, "F") +
            " - " +
            time(user.joinedAt, "R") +
            "\n**Roles:** " +
            rolesString,
        thumbnail: {
            url: "https://cdn.discordapp.com/emojis/1064442673806704672.webp"
        },
        color: 0x5865f2,
        author: {
            name: user.username,
            icon_url: user.avatarURL
        },
        footer: {
            text: "Event ID: " + event.id + " | " + eventType + " event"
        },
        timestamp: new Date()
    };

    const send = await sendMessage(
        {
            embeds: [embed]
        },
        settings.types.members.webhook_url
    ).catch((error) => {
        logger.error("Error sending " + eventType + " webhook", error);
        return { error: "Error sending " + eventType + " webhook" };
    });

    if (!send) {
        return { error: "Error sending " + eventType + " webhook" };
    }

    return true;
}

// member update event
async function memberUpdateEvent(data: any) {
    const eventType = "member_update";

    const guild = data.guild,
        oldMember = data.oldMember,
        newMember = data.newMember,
        newRoles = data.newRoles,
        oldRoles = data.oldRoles;

    let settings: any = await getMemberEventLogSettings(eventType, guild.id);
    if (!settings) {
        return;
    }

    settings = settings.settings;

    if (!settings.types.members) {
        return;
    }

    let descString = `**Member:** <@${newMember.userId}> (${newMember.userId})\n\`\`\`diff\n`;

    const embed = {
        title: "Member Updated",
        description: descString,
        thumbnail: {
            url: "https://cdn.discordapp.com/emojis/1064444245588578385.webp"
        },
        color: 0x4ca99d,
        author: {
            name: newMember.username,
            icon_url: newMember.avatarURL
        },
        footer: {
            text: "Event ID: " + newMember.userId + " | " + eventType + " event"
        },
        timestamp: new Date()
    };

    if (oldMember.nickname !== newMember.nickname) {
        embed.description += `\nNickname Changed:\n- ${oldMember.nickname ? oldMember.nickname : "None"}\n+ ${
            newMember.nickname ? newMember.nickname : "None"
        }\n`;
    }
    if (oldMember.avatar !== newMember.avatar) {
        embed.description += `\nAvatar Changed:\n- ${oldMember.avatarURL}\n+ ${newMember.avatarURL}\n`;
    }
    if (oldMember.username !== newMember.username) {
        embed.description += `\nUsername Changed:\n- ${oldMember.username}\n+ ${newMember.username}\n`;
    }

    embed.description += "```";

    if (oldRoles.length !== newRoles.length) {
        const addedRoles = newRoles
            .filter((role: any) => !oldRoles.includes(role))
            .map((role: any) => `<@&${role}>`)
            .join(", ");
        const removedRoles = oldRoles
            .filter((role: any) => !newRoles.includes(role))
            .map((role: any) => `<@&${role}>`)
            .join(", ");
        embed.description += `\n**Roles Changed:**`;
        if (addedRoles) {
            embed.description += `\n**Added:** ${addedRoles}`;
        }
        if (removedRoles) {
            embed.description += `\n**Removed:** ${removedRoles}`;
        }
    }

    const send = await sendMessage(
        {
            embeds: [embed]
        },
        settings.types.members.webhook_url
    ).catch((error) => {
        logger.error("Error sending " + eventType + " webhook", error);
        return { error: "Error sending " + eventType + " webhook" };
    });

    if (!send) {
        return { error: "Error sending " + eventType + " webhook" };
    }

    return true;
}

// export the events
export { memberBanEvent, memberUnbanEvent, memberJoinEvent, memberLeaveEvent, memberUpdateEvent };
