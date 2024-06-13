import winston from "winston";
import { WinstonTransport as AxiomTransport } from "@axiomhq/winston";
import config from "../config";

const logger = winston.createLogger({
    level: "info",
    format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.json()
        // winston.format.colorize({
        //     message: true
        // })
    ),
    defaultMeta: { service: "events-api" },
    transports: [
        new winston.transports.Console({ format: winston.format.simple() }),
        new AxiomTransport({
            dataset: config.axiom.dataset,
            token: config.axiom.token,
            orgId: config.axiom.orgId
        })
    ]
});

export default logger;
