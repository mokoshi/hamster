const { PHASE_DEVELOPMENT_SERVER, PHASE_PRODUCTION_BUILD } = require('next/constants');

module.exports = (phase) => {
  const isDev = phase === PHASE_DEVELOPMENT_SERVER;
  const isProd = phase === PHASE_PRODUCTION_BUILD && process.env.STAGING !== '1';
  const isStaging = phase === PHASE_PRODUCTION_BUILD && process.env.STAGING === '1';

  return {
    reactStrictMode: true,
    env: {
      // TODO 環境によって切り替える
      ApiBaseUrl: isDev ? 'http://localhost:4100' : 'http://localhost:4100',
    },
  };
};
