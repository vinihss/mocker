import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import CreateMock from './pages/CreateMock'
import MockDetail from './pages/MockDetail'
import MockTest from './pages/MockTest'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/mocks/new" element={<CreateMock />} />
        <Route path="/mocks/:id" element={<MockDetail />} />
        <Route path="/mocks/:id/test" element={<MockTest />} />
        <Route path="/mocks/:id/edit" element={<CreateMock isEdit />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App