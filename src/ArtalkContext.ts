import Artalk from './Artalk'

export default class ArtalkContext {
  public artalk: Artalk

  public constructor (artalk: Artalk) {
    this.artalk = artalk
  }
}
