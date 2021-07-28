const yaml = require('js-yaml')

export function configParse(configSourceCode) {
  try {
    const doc = yaml.load(configSourceCode);
    if (doc['version']) {
      return doc
    } 
    return null
  } catch (e) {
    return null
  }
}