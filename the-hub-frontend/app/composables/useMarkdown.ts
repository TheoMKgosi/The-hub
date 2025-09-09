import { marked } from 'marked'
import { mangle } from 'marked-mangle'
import { gfmHeadingId } from "marked-gfm-heading-id";
import { useMath } from './useMath'

export const useMarkdown = () => {
  const { renderMath } = useMath()

  const renderMarkdown = (text: string): string => {
    if (!text) return ''

    // Configure marked with custom renderer for better line break handling
    const renderer = new marked.Renderer()

    // Override paragraph rendering to preserve single line breaks
    renderer.paragraph = ({ tokens }) => {

      const text = tokens.map(token => token.raw).join('')
      // Replace single line breaks with <br> tags, but preserve double line breaks
      const processed = text.replace(/\n(?!\n)/g, '<br>')
      return `<p>${processed}</p>`
    }

    marked.setOptions({
      breaks: true,
      gfm: true,
      renderer: renderer
    })

    marked.use(mangle(), gfmHeadingId())

    // First render math formulas
    const textWithMath = renderMath(text)

    // Then render markdown
    return marked.parse(textWithMath) as string
  }

  return {
    renderMarkdown
  }
}
