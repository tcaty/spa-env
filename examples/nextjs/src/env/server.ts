import zod from 'zod'

const schema = zod.object({
    POSTGRES_CONN_STRING: zod.string(),
}).readonly()

export const env = schema.parse({
    POSTGRES_CONN_STRING: process.env.POSTGRES_CONN_STRING
})
