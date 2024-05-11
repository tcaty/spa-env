'use client'

import { env } from '@/env/client'
import React from 'react'

export const ClientSideComponent = () => {
	return (
		<div>
            <h2>client side component</h2>
            <h3>client side variables</h3>
            <h4>{env.API_URL}</h4>
            <h4>{env.SECRET_TOKEN}</h4>
        </div>
	)
}