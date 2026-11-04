import { describe, it, expect } from 'vitest'
import { useMarkdown } from '../../../app/composables/useMarkdown'

describe('useMarkdown', () => {
  it('returns an empty string for empty input', () => {
    const { renderMarkdown } = useMarkdown()

    expect(renderMarkdown('')).toBe('')
  })

  it('wraps plain text in a paragraph', () => {
    const { renderMarkdown } = useMarkdown()

    expect(renderMarkdown('Hello world')).toBe('<p>Hello world</p>')
  })

  it('converts single line breaks to <br> within a paragraph', () => {
    const { renderMarkdown } = useMarkdown()

    expect(renderMarkdown('line one\nline two')).toBe('<p>line one<br>line two</p>')
  })

  it('renders headings with ids from the gfm plugin', () => {
    const { renderMarkdown } = useMarkdown()

    const html = renderMarkdown('# My Heading')

    expect(html).toContain('<h1')
    expect(html).toContain('id=')
    expect(html).toContain('My Heading</h1>')
  })

  // NOTE: current behavior — the custom paragraph renderer joins raw tokens,
  // so inline elements (code, emphasis) are NOT parsed inside paragraphs.
  // This documents actual output; see flagged bug if inline parsing is desired.
  it('keeps inline code literal inside paragraphs (raw token join)', () => {
    const { renderMarkdown } = useMarkdown()

    expect(renderMarkdown('Run `npm test` now')).toBe('<p>Run `npm test` now</p>')
  })

  it('renders inline math expressions via katex', () => {
    const { renderMarkdown } = useMarkdown()

    const html = renderMarkdown('Euler: $e^{i\\pi}$')

    expect(html).toContain('katex')
  })

  it('keeps emphasis literal inside paragraphs (raw token join)', () => {
    const { renderMarkdown } = useMarkdown()

    expect(renderMarkdown('this is **important**')).toBe('<p>this is **important**</p>')
  })
})
