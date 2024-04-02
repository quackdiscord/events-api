const config = {
    port: process.env.PORT || 3000,
    token: process.env.TOKEN,

    db: {
        host: process.env.DATABASE_HOST,
        user: process.env.DATABASE_USERNAME,
        database: process.env.DATABASE_NAME,
        password: process.env.DATABASE_PASSWORD,
        uri: process.env.DATABASE_URI
    },

    redis: {
        url: process.env.REDIS_HOST,
        token: process.env.REDIS_TOKEN
    },

    kafka: {
        brokers: [process.env.KAFKA_BROKER as string],
        username: process.env.KAFKA_USERNAME as string,
        password: process.env.KAFKA_PASSWORD as string
    },

    axiom: {
        dataset: process.env.AXIOM_DATASET,
        token: process.env.AXIOM_TOKEN as string,
        orgId: process.env.AXIOM_ORGID
    }
};

export default config;
