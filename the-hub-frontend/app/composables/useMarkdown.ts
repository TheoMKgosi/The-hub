import { marked } from 'marked'
import { useMath } from './useMath'

export const useMarkdown = () => {
  const { renderMath } = useMath()

  const renderMarkdown = (text: string): string => {
    if (!text) return ''

    // Configure marked with custom renderer for better line break handling
    const renderer = new marked.Renderer()

    // Override paragraph rendering to preserve single line breaks
    renderer.paragraph = (text: string) => {
      // Replace single line breaks with <br> tags, but preserve double line breaks
      const processed = text.replace(/\n(?!\n)/g, '<br>')
      return `<p>${processed}</p>`
    }

    marked.setOptions({
      breaks: true,
      gfm: true,
      headerIds: false,
      mangle: false,
      renderer: renderer
    })

    // First render math formulas
    const textWithMath = renderMath(text)

    // Then render markdown
    return marked.parse(textWithMath) as string
  }

  return {
    renderMarkdown
  }
}