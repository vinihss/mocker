import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, Play, Copy, Check, Clock } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { mocksApi } from '@/lib/api'
import type { Mock, TestMockResponse } from '@/types/api'

export default function MockTest() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [mock, setMock] = useState<Mock | null>(null)
  const [loading, setLoading] = useState(true)
  const [testing, setTesting] = useState(false)
  const [input, setInput] = useState('{}')
  const [response, setResponse] = useState<TestMockResponse | null>(null)
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
      } finally {
        setLoading(false)
      }
    }
    fetchMock()
  }, [id])

  const handleTest = async () => {
    if (!mock) return
    setTesting(true)
    setError(null)
    setResponse(null)

    try {
      const inputObj = JSON.parse(input)
      const result = await mocksApi.test(mock.id, { input: inputObj })
      setResponse(result)
    } catch (err) {
      console.error('Failed to test mock:', err)
      setError(err instanceof Error ? err.message : 'Failed to test mock')
    } finally {
      setTesting(false)
    }
  }

  const handleCopyResponse = () => {
    if (!response) return
    navigator.clipboard.writeText(response.output)
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

  if (!mock) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card>
          <CardContent className="pt-6">
            <p className="text-destructive mb-4">Mock not found</p>
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
          <Button variant="ghost" onClick={() => navigate(`/mocks/${id}`)}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
        </div>
      </header>

      <main className="container mx-auto px-4 py-6">
        <div className="flex items-center gap-3 mb-6">
          <Badge variant={mock.method === 'GET' ? 'default' : 'secondary'} className="font-mono">
            {mock.method}
          </Badge>
          <h1 className="text-2xl font-bold">{mock.name}</h1>
          <code className="bg-muted px-2 py-1 rounded text-sm">{mock.path}</code>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Request Input</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <p className="text-sm text-muted-foreground">
                  JSON input to send to the mock endpoint
                </p>
                <textarea
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  className="font-mono min-h-[300px] w-full"
                  placeholder='{"key": "value"}'
                />
              </div>
              <Button onClick={handleTest} disabled={testing} className="w-full">
                <Play className="h-4 w-4 mr-2" />
                {testing ? 'Testing...' : 'Send Request'}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>Response</CardTitle>
              {response && (
                <div className="flex items-center gap-2">
                  <Badge
                    variant={response.statusCode < 400 ? 'success' : 'destructive'}
                  >
                    {response.statusCode}
                  </Badge>
                  <span className="text-sm text-muted-foreground flex items-center">
                    <Clock className="h-3 w-3 mr-1" />
                    {response.tookMs}ms
                  </span>
                </div>
              )}
            </CardHeader>
            <CardContent>
              {error ? (
                <div className="p-3 rounded-md bg-destructive/10 text-destructive text-sm">
                  {error}
                </div>
              ) : response ? (
                <div className="space-y-4">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleCopyResponse}
                    className="mb-2"
                  >
                    {copied ? (
                      <Check className="h-4 w-4 mr-2" />
                    ) : (
                      <Copy className="h-4 w-4 mr-2" />
                    )}
                    {copied ? 'Copied!' : 'Copy'}
                  </Button>
                  <pre className="bg-muted p-4 rounded-md overflow-x-auto text-sm font-mono max-h-[400px] overflow-y-auto">
                    {JSON.stringify(JSON.parse(response.output), null, 2)}
                  </pre>
                </div>
              ) : (
                <p className="text-muted-foreground text-center py-12">
                  Send a request to see the response
                </p>
              )}
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  )
}