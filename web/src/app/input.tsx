"use client"
import React from 'react'

function Input({ onChange} : { onChange: (e: React.ChangeEvent<HTMLInputElement>) => void }){
  return (
    <div>
        <input type="file" onChange={onChange}/>
    </div>
  )
}

export default Input