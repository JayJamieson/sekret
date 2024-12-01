let onCopy = () => {}

const copy = async (target, mimeType = undefined) => {
    if (typeof target === 'function') {
        target = target()
    }

    if (typeof target === 'object') {
        target = JSON.stringify(target)
    }

    if (mimeType !== undefined) {
        const result = await window.navigator.clipboard.write([
        new ClipboardItem({
          [mimeType]: new Blob([target], {
            type: mimeType,
          })
        })
      ])
      return onCopy(result)
    }

    result = await window.navigator.clipboard.writeText(target)
  return onCopy(result)
}

function Clipboard(Alpine) {
    Alpine.magic('clipboard', () => {
        return copy
    })

    Alpine.directive('clipboard', (el, { modifiers, expression }, { evaluateLater, cleanup }) => {
        const getCopyContent = modifiers.includes('raw') ? c => c(expression) : evaluateLater(expression)
        const clickHandler = () => getCopyContent(copy)

        el.addEventListener('click', clickHandler)

        cleanup(() => {
            el.removeEventListener('click', clickHandler)
        })
    })
}

Clipboard.configure = (config) => {
    if (config.hasOwnProperty('onCopy') && typeof config.onCopy === 'function') {
        onCopy = config.onCopy
    }

    return Clipboard
}

export default Clipboard;
