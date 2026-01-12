// Add to document: main.tsx setup example
import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createBrowserRouter } from 'react-router'
import { generatedRoutes } from './router'
import './global.css'

const router = createBrowserRouter(generatedRoutes)

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)