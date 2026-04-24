import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Plus, RefreshCw, Power, PowerOff, Trash2, TestTube } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { mocksApi } from '@/lib/api'
import type { Mock } from '@/types/api'

export default function Dashboard() {
  const [mocks, setMocks] = useState<Mock[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState<string>('all')
  const [error, setError] = useState<string | null>(null)

  const fetchMocks = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await mocksApi.list()
      setMocks(data)
    } catch (err) {
      setError('Failed to load mocks. Make sure the backend is running on port 8080.')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchMocks()
  }, [])

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this mock?')) return
    try {
      await mocksApi.delete(id)
      setMocks(mocks.filter((m) => m.id !== id))
    } catch (err) {
      console.error('Failed to delete mock:', err)
    }
  }

  const handleToggleActive = async (mock: Mock) => {
    try {
      const updated = mock.isActive
        ? await mocksApi.deactivate(mock.id)
        : await mocksApi.activate(mock.id)
      setMocks(mocks.map((m) => (m.id === updated.id ? updated : m)))
    } catch (err) {
      console.error('Failed to toggle mock:', err)
    }
  }

  const filteredMocks = mocks.filter((mock) => {
    const matchesSearch =
      mock.name.toLowerCase().includes(search.toLowerCase()) ||
      mock.path.toLowerCase().includes(search.toLowerCase())
    const matchesFilter =
      filter === 'all' ||
      (filter === 'active' && mock.isActive) ||
      (filter === 'inactive' && !mock.isActive)
    return matchesSearch && matchesFilter
  })

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b bg-card">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold">Mock API Dashboard</h1>
            <p className="text-sm text-muted-foreground">
              Gerencie seus mocks dinâmicos
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={fetchMocks} disabled={loading}>
              <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
            <Link to="/mocks/new">
              <Button size="sm">
                <Plus className="h-4 w-4" />
                New Mock
              </Button>
            </Link>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-6">
        <div className="flex gap-4 mb-6">
          <div className="flex-1">
            <Input
              placeholder="Search mocks..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
          <Select value={filter} onValueChange={setFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All</SelectItem>
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="inactive">Inactive</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {error && (
          <Card className="mb-6 border-destructive">
            <CardContent className="pt-6">
              <p className="text-destructive">{error}</p>
              <Button variant="outline" className="mt-2" onClick={fetchMocks}>
                Try Again
              </Button>
            </CardContent>
          </Card>
        )}

        {loading ? (
          <div className="text-center py-12">
            <RefreshCw className="h-8 w-8 animate-spin mx-auto text-muted-foreground" />
            <p className="mt-2 text-muted-foreground">Loading...</p>
          </div>
        ) : filteredMocks.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <TestTube className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
              <p className="text-muted-foreground mb-4">No mocks found</p>
              <Link to="/mocks/new">
                <Button>Create your first mock</Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {filteredMocks.map((mock) => (
              <Card key={mock.id}>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div>
                        <div className="flex items-center gap-2">
                          <Badge
                            variant={mock.method === 'GET' ? 'default' : 'secondary'}
                            className="font-mono"
                          >
                            {mock.method}
                          </Badge>
                          <span className="font-medium">{mock.name}</span>
                          <Badge
                            variant={mock.isActive ? 'success' : 'outline'}
                          >
                            {mock.isActive ? 'Active' : 'Inactive'}
                          </Badge>
                        </div>
                        <p className="text-sm text-muted-foreground font-mono mt-1">
                          {mock.path}
                        </p>
                        {mock.description && (
                          <p className="text-sm text-muted-foreground mt-1">
                            {mock.description}
                          </p>
                        )}
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Link to={`/mocks/${mock.id}/test`}>
                        <Button variant="outline" size="sm">
                          <TestTube className="h-4 w-4" />
                        </Button>
                      </Link>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleToggleActive(mock)}
                      >
                        {mock.isActive ? (
                          <PowerOff className="h-4 w-4" />
                        ) : (
                          <Power className="h-4 w-4" />
                        )}
                      </Button>
                      <Link to={`/mocks/${mock.id}`}>
                        <Button variant="outline" size="sm">View</Button>
                      </Link>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDelete(mock.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}