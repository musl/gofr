import Axios from 'axios'
import Ractive from "ractive"
import Router from "ractive-route"

import Heading from "./Heading"
import Sidebar from "./Sidebar"
import Character from "./Character"
import Inventory from "./Inventory"

import template from "./App.html"

import "./pure.css"
import "./grids-responsive.css"
import "./App.css"

export default Ractive.extend({
  template: template,
  components: {
    Heading,
    Sidebar,
  },
  on: {
    complete: function(context) {
      var self = this;
      
      Axios.get('/api/data')
        .then(function(response) {
          if(response.data.name) {
            self.set(response.data);
            self.route();
          } else {
            console.log("No data, response: ", response);
          }
        })
        .catch(function(error) {
          console.log(error);
        });
    },
  },
  route: function() {
    const router = new Router({
      el: "#router",
      basepath: "/",
      data: this.get(),
    });

    router.addRoute('/?', Character);
    router.addRoute('/inventory', Inventory);

    router.init({
      noHistory: true,
      reload: false,
    });

    router.watchLinks();
    router.watchState();
  },
});
