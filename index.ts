import log from "./lib/logger";
import { kafka } from "./lib/kafka";

// import the events
import * as ServerEvents from "./events/serverEvents";
import * as MemberEvents from "./events/memberEvents";
import * as MessageEvents from "./events/messageEvents";

// connect to kafka
const kfReceiver = kafka.consumer({ groupId: "consumer-1" });

await kfReceiver.connect();

// reciever events
kfReceiver.on("consumer.connect", () => {
    log.info(`Consumer connected`);
});

kfReceiver.on("consumer.crash", (error) => {
    log.error(`Consumer crashed`, error);
});

kfReceiver.on("consumer.disconnect", (error) => {
    log.error(`Consumer disconnected`, error);
});

kfReceiver.subscribe({ topic: "event-logs" });

// for each message we get, run the event then set a 1 second timeout
await kfReceiver.run({
    eachMessage: async ({ message }: any) => {
        try {
            const data = JSON.parse(message?.value?.toString() as string);

            log.info(`Received event: ${data.type}`);

            switch (data.type) {
                case "server_update":
                    ServerEvents.serverUpdateEvent(data);
                    break;
                case "channel_create":
                    ServerEvents.channelCreateEvent(data);
                    break;
                case "channel_delete":
                    ServerEvents.channelDeleteEvent(data);
                    break;
                case "channel_pins_update":
                    ServerEvents.channelPinsUpdateEvent(data);
                    break;
                case "channel_update":
                    ServerEvents.channelUpdateEvent(data);
                    break;
                case "emoji_create":
                    ServerEvents.emojiCreateEvent(data);
                    break;
                case "emoji_delete":
                    ServerEvents.emojiDeleteEvent(data);
                    break;
                case "emoji_update":
                    ServerEvents.emojiUpdateEvent(data);
                    break;
                case "event_create":
                    ServerEvents.eventCreateEvent(data);
                    break;
                case "event_delete":
                    ServerEvents.eventDeleteEvent(data);
                    break;
                case "invite_create":
                    ServerEvents.inviteCreateEvent(data);
                    break;
                case "invite_delete":
                    ServerEvents.inviteDeleteEvent(data);
                    break;
                case "role_create":
                    ServerEvents.roleCreateEvent(data);
                    break;
                case "role_delete":
                    ServerEvents.roleDeleteEvent(data);
                    break;
                case "role_update":
                    ServerEvents.roleUpdateEvent(data);
                    break;
                case "member_ban":
                    MemberEvents.memberBanEvent(data);
                    break;
                case "member_unban":
                    MemberEvents.memberUnbanEvent(data);
                    break;
                case "member_join":
                    MemberEvents.memberJoinEvent(data);
                    break;
                case "member_leave":
                    MemberEvents.memberLeaveEvent(data);
                    break;
                case "member_update":
                    MemberEvents.memberUpdateEvent(data);
                    break;
                case "message_delete":
                    MessageEvents.messageDeleteEvent(data);
                    break;
                case "message_update":
                    MessageEvents.messageUpdateEvent(data);
                    break;
                case "message_bulk_delete":
                    MessageEvents.messageBulkDeleteEvent(data);
                    break;
                default:
                    log.warn(`Received unknown event type: ${data.type}`);
                    break;
            }
        } catch (error) {
            log.error(`Error parsing event`, error);
        }
    }
});
