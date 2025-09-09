import katex from 'katex'
import 'katex/dist/katex.min.css'

export const useMath = () => {
  const renderMath = (text: string): string => {
    if (!text) return ''

    // Replace $...$ with rendered KaTeX
    const mathRegex = /\$([^$]+)\$/g
    let result = text
    let match

    while ((match = mathRegex.exec(text)) !== null) {
      try {
        const mathHtml = katex.renderToString(match[1], {
          throwOnError: false,
          displayMode: false
        })
        result = result.replace(match[0], mathHtml)
      } catch (error) {
        console.warn('Math rendering error:', error)
        // Keep original text if rendering fails
      }
    }

    return result
  }

  return {
    renderMath
  }
}