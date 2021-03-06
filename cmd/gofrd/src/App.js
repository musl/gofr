import Ractive from "ractive"

import Heading from "./Heading"
import Browser from "./Browser"

import "./pure.css"
import "./grids-responsive.css"

import template from "./App.html"
import "./App.css"

export default Ractive.extend({
	template: template,
	components: {
		Heading,
		Browser,
	},
	hello: {},
	on: {},
});
