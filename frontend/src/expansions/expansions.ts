import {EXPANSION_ICONS_SMALL} from "./eq-expansion-icons";
import {EXPANSION_NAMES} from "./eq-expansions";

export default class Expansions {
  static getExpansionIconUrlSmall(expansionId: string | number) {
    // @ts-ignore
    if (EXPANSION_ICONS_SMALL[expansionId]) {
      // @ts-ignore
      return require('@/assets/images/expansions/' + EXPANSION_ICONS_SMALL[expansionId])
    }

    // return transparent base64 encoded image if nothing found
    return 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7'
  }
  static getExpansionName(expansionId: number) {
    if (EXPANSION_NAMES[expansionId]) {
      return EXPANSION_NAMES[expansionId]
    }

    if (expansionId === -1) {
      return 'All'
    }

    // return unknown expansion if not found
    return 'Unknown'
  }
}
