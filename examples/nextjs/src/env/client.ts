import zod from 'zod'

const schema = zod.object({
    API_URL: zod.string(),
    SECRET_TOKEN: zod.string(),
}).readonly()

export const env = schema.parse({
    API_URL: `"${process.env.NEXT_PUBLIC_API_URL}"`,
    SECRET_TOKEN: process.env.NEXT_PUBLIC_SECRET_TOKEN,
})
