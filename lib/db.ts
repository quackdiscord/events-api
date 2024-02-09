import { drizzle } from "drizzle-orm/planetscale-serverless";
import { connect } from "@planetscale/database";
import config from "../config";

// create the connection
const connection = connect({
    host: config.database.host,
    username: config.database.username,
    password: config.database.password
});

const db = drizzle(connection);

export { db };
