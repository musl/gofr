import Ractive from 'ractive'

import App from './App'

const app = App({ el: '#app' })

if(module.hot) {
  module.hot.accept();
}

