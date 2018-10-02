class Artalk {
  constructor (opts) {
    this.opts = opts
    this.init()
  }

  init () {
    this.el = ({}).toString.call(this.opts) === '[object HTMLDivElement]' ? this.opts.el : document.querySelectorAll(this.opts.el)[0]
    if (({}).toString.call(this.el) !== '[object HTMLDivElement]') {
      throw Error(`Sorry, Target element "${this.opts.el}" was not found.`)
    }
    this.el.className += 'artalk'
    this.el.innerHTML = 'Hello World'
    console.log('Hello World')
  }
}

module.exports = Artalk
