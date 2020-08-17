
export const Config = JSON.stringify(process.env.NODE_ENV === 'production' ? require('../config.prod.json') : require('../config.dev.json'))
