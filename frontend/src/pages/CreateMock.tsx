import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Plus, Trash2, Save } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { mocksApi } from '@/lib/api'
import type { ResponseConfig, FakerField } from '@/types/api'

interface CreateMockProps {
  isEdit?: boolean
}

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'] as const
const RESPONSE_TYPES = ['static', 'dynamic', 'faker'] as const
const FAKER_TYPES = [
  'name', 'first_name', 'last_name', 'email', 'phone', 'phone_number',
  'uuid', 'uuid_rfc4122', 'date', 'datetime', 'timestamp', 'word',
  'sentence', 'paragraph', 'address', 'city', 'country', 'state',
  'zipcode', 'latitude', 'longitude', 'url', 'ip', 'ipv4', 'ipv6',
  'credit_card', 'credit_card_number', 'credit_card_type'
] as const

export default function CreateMock({ isEdit }: CreateMockProps) {
  const { id } = useParams()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(isEdit)
  const [error, setError] = useState<string | null>(null)

  const [name, setName] = useState('')
  const [path, setPath] = useState('')
  const [method, setMethod] = useState<string>('GET')
  const [description, setDescription] = useState('')

  const [responseType, setResponseType] = useState<ResponseConfig['type']>('static')
  const [statusCode, setStatusCode] = useState(200)
  const [delayMs, setDelayMs] = useState(0)
  const [responseBody, setResponseBody] = useState('')
  const [fakerFields, setFakerFields] = useState<FakerField[]>([])

  useEffect(() => {
    if (!isEdit || !id) return

    const fetchMock = async () => {
      try {
        const mock = await mocksApi.get(id)
        setName(mock.name)
        setPath(mock.path)
        setMethod(mock.method)
        setDescription(mock.description || '')
        setResponseType(mock.responseConfig.type)
        setStatusCode(mock.responseConfig.statusCode)
        setDelayMs(mock.responseConfig.delayMs || 0)
        setResponseBody(JSON.stringify(mock.responseConfig.body, null, 2))
        if (mock.responseConfig.faker?.fields) {
          setFakerFields(mock.responseConfig.faker.fields)
        }
      } catch (err) {
        console.error('Failed to fetch mock:', err)
        setError('Failed to load mock')
      } finally {
        setFetching(false)
      }
    }

    fetchMock()
  }, [isEdit, id])

  const handleAddFakerField = () => {
    setFakerFields([...fakerFields, { name: '', type: 'name' }])
  }

  const handleRemoveFakerField = (index: number) => {
    setFakerFields(fakerFields.filter((_, i) => i !== index))
  }

  const handleFakerFieldChange = (index: number, field: keyof FakerField, value: string) => {
    const updated = [...fakerFields]
    updated[index] = { ...updated[index], [field]: value }
    setFakerFields(updated)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      let body: unknown
      if (responseType === 'faker') {
        body = { faker: { fields: fakerFields } }
      } else {
        body = JSON.parse(responseBody)
      }

      const responseConfig: ResponseConfig = {
        type: responseType,
        statusCode,
        delayMs: delayMs || undefined,
        body,
      }

      if (isEdit && id) {
        await mocksApi.update(id, {
          name,
          path,
          method,
          description,
          responseConfig,
        })
      } else {
        await mocksApi.create({
          name,
          path,
          method,
          description,
          responseConfig,
        })
      }

      navigate('/')
    } catch (err) {
      console.error('Failed to save mock:', err)
      setError(err instanceof Error ? err.message : 'Failed to save mock')
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p>Loading...</p>
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

      <main className="container mx-auto px-4 py-6 max-w-3xl">
        <form onSubmit={handleSubmit}>
          <Card>
            <CardHeader>
              <CardTitle>{isEdit ? 'Edit Mock' : 'Create New Mock'}</CardTitle>
              <CardDescription>
                Configure o endpoint do mock com resposta estática, dinâmica ou dados fake
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {error && (
                <div className="p-3 rounded-md bg-destructive/10 text-destructive text-sm">
                  {error}
                </div>
              )}

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input
                    id="name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="get-users"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Input
                    id="description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="Optional description"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="method">Method</Label>
                  <Select value={method} onValueChange={setMethod}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {HTTP_METHODS.map((m) => (
                        <SelectItem key={m} value={m}>{m}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="path">Path</Label>
                  <Input
                    id="path"
                    value={path}
                    onChange={(e) => setPath(e.target.value)}
                    placeholder="/api/users"
                    required
                  />
                </div>
              </div>

              <div className="border-t pt-6">
                <h3 className="text-lg font-medium mb-4">Response Config</h3>

                <div className="grid grid-cols-3 gap-4 mb-4">
                  <div className="space-y-2">
                    <Label>Type</Label>
                    <Select value={responseType} onValueChange={(v) => setResponseType(v as ResponseConfig['type'])}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {RESPONSE_TYPES.map((t) => (
                          <SelectItem key={t} value={t}>
                            {t.charAt(0).toUpperCase() + t.slice(1)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>Status Code</Label>
                    <Input
                      type="number"
                      value={statusCode}
                      onChange={(e) => setStatusCode(Number(e.target.value))}
                      min={100}
                      max={599}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Delay (ms)</Label>
                    <Input
                      type="number"
                      value={delayMs}
                      onChange={(e) => setDelayMs(Number(e.target.value))}
                      min={0}
                    />
                  </div>
                </div>

                {responseType === 'faker' ? (
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <Label>Faker Fields</Label>
                      <Button type="button" variant="outline" size="sm" onClick={handleAddFakerField}>
                        <Plus className="h-4 w-4 mr-2" />
                        Add Field
                      </Button>
                    </div>
                    {fakerFields.length === 0 ? (
                      <p className="text-sm text-muted-foreground">No fields added yet</p>
                    ) : (
                      <div className="space-y-2">
                        {fakerFields.map((field, index) => (
                          <div key={index} className="flex gap-2">
                            <Input
                              placeholder="Field name"
                              value={field.name}
                              onChange={(e) => handleFakerFieldChange(index, 'name', e.target.value)}
                            />
                            <Select value={field.type} onValueChange={(v) => handleFakerFieldChange(index, 'type', v)}>
                              <SelectTrigger className="w-[180px]">
                                <SelectValue />
                              </SelectTrigger>
                              <SelectContent>
                                {FAKER_TYPES.map((t) => (
                                  <SelectItem key={t} value={t}>{t}</SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              onClick={() => handleRemoveFakerField(index)}
                            >
                              <Trash2 className="h-4 w-4 text-destructive" />
                            </Button>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="space-y-2">
                    <Label>Response Body (JSON)</Label>
                    <Textarea
                      value={responseBody}
                      onChange={(e) => setResponseBody(e.target.value)}
                      placeholder='{"message": "Hello World"}'
                      className="font-mono min-h-[200px]"
                      required
                    />
                  </div>
                )}
              </div>

              <div className="flex justify-end gap-2 pt-4 border-t">
                <Button type="button" variant="outline" onClick={() => navigate('/')}>
                  Cancel
                </Button>
                <Button type="submit" disabled={loading}>
                  <Save className="h-4 w-4 mr-2" />
                  {loading ? 'Saving...' : 'Save Mock'}
                </Button>
              </div>
            </CardContent>
          </Card>
        </form>
      </main>
    </div>
  )
}