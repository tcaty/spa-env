import { env } from '@/env/server'
import React from 'react'

export const ServerSideComponent = () => {
	return (
		<div>
            <h2>server side component</h2>
            <h3>server side variables</h3>
            <h4>{env.POSTGRES_CONN_STRING}</h4>
        </div>
	)
}