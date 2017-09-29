import Ractive from 'ractive'

import App from './App'
import Defaults from './Defaults'

const app = App({
  el: '#app',
  data: Defaults,
});

if(module.hot) {
  module.hot.accept();
}

