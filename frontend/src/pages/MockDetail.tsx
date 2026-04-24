import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { ArrowLeft, Edit, Power, PowerOff, Trash2, TestTube, Copy, Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { mocksApi } from '@/lib/api'
import type { Mock } from '@/types/api'

export default function MockDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [mock, setMock] = useState<Mock | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    const fetchMock = async () => {
      if (!id) return
      try {
        const data = await mocksApi.get(id)
        setMock(data)
      } catch (err) {
        console.error('Failed to fetch mock:', err)
        setError('Mock not found')
      } finally {
        setLoading(false)
      }
    }
    fetchMock()
  }, [id])

  const handleDelete = async () => {
    if (!mock || !confirm('Are you sure you want to delete this mock?')) return
    try {
      await mocksApi.delete(mock.id)
      navigate('/')
    } catch (err) {
      console.error('Failed to delete mock:', err)
    }
  }

  const handleToggleActive = async () => {
    if (!mock) return
    try {
      const updated = mock.isActive
        ? await mocksApi.deactivate(mock.id)
        : await mocksApi.activate(mock.id)
      setMock(updated)
    } catch (err) {
      console.error('Failed to toggle mock:', err)
    }
  }

  const handleCopy = () => {
    if (!mock) return
    navigator.clipboard.writeText(`${mock.method} ${mock.path}`)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p>Loading...</p>
      </div>
    )
  }

  if (error || !mock) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card>
          <CardContent className="pt-6">
            <p className="text-destructive mb-4">{error || 'Mock not found'}</p>
            <Button onClick={() => navigate('/')}>Back to Dashboard</Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b bg-card">
        <div className="container mx-auto px-4 py-4">
          <Button variant="ghost" onClick={() => navigate('/')}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
        </div>
      </header>

      <main className="container mx-auto px-4 py-6">
        <div className="flex items-start justify-between mb-6">
          <div>
            <div className="flex items-center gap-3">
              <Badge variant={mock.method === 'GET' ? 'default' : 'secondary'} className="font-mono text-lg">
                {mock.method}
              </Badge>
              <h1 className="text-2xl font-bold">{mock.name}</h1>
              <Badge variant={mock.isActive ? 'success' : 'outline'}>
                {mock.isActive ? 'Active' : 'Inactive'}
              </Badge>
            </div>
            <div className="flex items-center gap-2 mt-2">
              <code className="bg-muted px-2 py-1 rounded text-sm">{mock.path}</code>
              <Button variant="ghost" size="sm" onClick={handleCopy}>
                {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
              </Button>
            </div>
            {mock.description && (
              <p className="text-muted-foreground mt-2">{mock.description}</p>
            )}
          </div>
          <div className="flex items-center gap-2">
            <Link to={`/mocks/${mock.id}/test`}>
              <Button variant="outline">
                <TestTube className="h-4 w-4 mr-2" />
                Test
              </Button>
            </Link>
            <Button variant="outline" onClick={handleToggleActive}>
              {mock.isActive ? (
                <PowerOff className="h-4 w-4 mr-2" />
              ) : (
                <Power className="h-4 w-4 mr-2" />
              )}
              {mock.isActive ? 'Deactivate' : 'Activate'}
            </Button>
            <Link to={`/mocks/${mock.id}/edit`}>
              <Button variant="outline">
                <Edit className="h-4 w-4 mr-2" />
                Edit
              </Button>
            </Link>
            <Button variant="destructive" onClick={handleDelete}>
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>

        <Tabs defaultValue="config">
          <TabsList>
            <TabsTrigger value="config">Config</TabsTrigger>
            <TabsTrigger value="response">Response</TabsTrigger>
          </TabsList>

          <TabsContent value="config">
            <Card>
              <CardHeader>
                <CardTitle>Configuration</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Created</p>
                    <p>{new Date(mock.createdAt).toLocaleString()}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Updated</p>
                    <p>{new Date(mock.updatedAt).toLocaleString()}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Method</p>
                    <p className="font-mono">{mock.method}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Path</p>
                    <p className="font-mono">{mock.path}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="response">
            <Card>
              <CardHeader>
                <CardTitle>Response Configuration</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-3 mb-4">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Type</p>
                    <Badge variant="secondary">{mock.responseConfig.type}</Badge>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Status Code</p>
                    <p className="font-mono">{mock.responseConfig.statusCode}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Delay</p>
                    <p>{mock.responseConfig.delayMs || 0}ms</p>
                  </div>
                </div>

                <div className="space-y-2">
                  <p className="text-sm font-medium text-muted-foreground">Body</p>
                  <pre className="bg-muted p-4 rounded-md overflow-x-auto text-sm font-mono">
                    {JSON.stringify(mock.responseConfig.body, null, 2)}
                  </pre>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </main>
    </div>
  )
}