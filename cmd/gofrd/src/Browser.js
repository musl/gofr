import Ractive from "ractive";

import Gofr from "./Gofr";
import Editor from "./Editor";

import template from "./Browser.html";

export default Ractive.extend({
	template: template,
	components: {
		Editor,
	},
	data: function() {
		return {
			width: 960,
			height: 960,
			view: {},
			default_view: {
				editable: false,
				i: 300,
				e: 4.0,
				m: '#444444',
				c: 'stripe',
				r: 'mandelbrot',
				s: 2.0,
				p: 2,
				rmin: -2.1,
				rmax: 0.6,
				imin: -1.25,
				imax: 1.25,
			},
			bookmarks: {},
			render_id: '',
			view_url: this.view_url,
		};
	},
	on: {
		render: function() {
			var marks, self, view, render_id;

			self = this;

			view = JSON.parse(Gofr.storage.getItem('gofr.browser.view'));
			if(view) {
				this.set('view', view);
			} else {
				this.copy_view('default_view', 'view');
			}

			marks = JSON.parse(Gofr.storage.getItem('gofr.browser.marks'));
			if(marks) {
				window.marks = marks;
				Object.keys(marks).filter(function(key, value) {
					self.set('bookmarks.' + key, value);
				});
			} else {
				this.copy_view('default_view', 'bookmarks.home');
			}

			render_id = JSON.parse(Gofr.storage.getItem('gofr.browser.render_id'));
			if(render_id) {
				this.set('render_id', render_id);
			} else {
				this.set('render_id', Gofr.uuid());
			}

			this.canvas = this.find('canvas');
			this.ctx = this.canvas.getContext('2d');

			this.ring = this.find('div#ring');
		},
		complete: function() {
			this.observe('view.*', function(newValue, oldValue, keypath) {

				// Filter out text fields so that we don't kick off extra render
				// jobs. Also filter out width and height because update_view()
				// changes them.
				//
				// TODO: Impose a delayed queue to coalesce edits or show a ui
				// element to indicate that the user needs to call for a refresh
				if(keypath.match(/\.(i|e|s|p|w|h)$/)) return;

				Gofr.storage.setItem('gofr.browser.view', this.json('view'));
				this.update_view();
			});

			this.observe('bookmarks', function() {
				Gofr.storage.setItem('gofr.browser.marks', this.json('bookmarks'));
			});

			this.observe('render_id', function() {
				Gofr.storage.setItem('gofr.browser.render_id', this.json('render_id'));
			});
		},
		move_up: function() {
			this.translate_view(0.0, -0.0625 * (this.get('view.imax') - this.get('view.imin')));
		},
		move_down: function() {
			this.translate_view(0.0, 0.0625 * (this.get('view.imax') - this.get('view.imin')));
		},
		move_left: function() {
			this.translate_view(-0.0625 * (this.get('view.rmax') - this.get('view.rmin')), 0.0);
		},
		move_right: function() {
			this.translate_view(0.0625 * (this.get('view.rmax') - this.get('view.rmin')), 0.0);
		},
		zoom_in: function() {
			this.scale_view(0.9);
		},
		zoom_in_4x: function() {
			this.scale_view(0.6);
		},
		zoom_out: function() {
			this.scale_view(1.1);
		},
		zoom_out_4x: function() {
			this.scale_view(1.4);
		},	
		update_view: function() {
			this.update_view();
		},
		go_to_bookmark: function(event) {
			var bookmark, name;

			name = 'bookmarks.' + event.node.data('bookmark');
			if(!name in this.get('bookmarks')) { return; }
			this.copy_view(name, 'view');
		},
		add_bookmark: function(event) {
			var name;

			// TODO FIX THIS window.prompt , turn it into a modal dialog.
			name = window.prompt('Enter a name:');

			if(name) {
				this.copy_view('view', 'bookmarks.' + name);
				this.set('bookmarks.' + name + '.editable', true); 
			}
		},
		update_bookmark: function(event) {
			var name;

			name = event.node.data('bookmark');
			this.copy_view('view', 'bookmarks.' + name);
			this.set('bookmarks.' + name + '.editable', true); 
		},
		delete_bookmark: function(event) {
			var bookmarks;

			bookmarks = this.get('bookmarks');
			delete bookmarks[event.node.data('bookmark')];
			this.update('bookmarks');
		},
		edit_view: function() {
			var editor;

			editor = this.findComponent('Editor');
			editor.fire('edit', 'view', this.json('view'));
		},
		edit_bookmarks: function() {
			var editor;

			editor = this.findComponent('editor');
			editor.fire('edit', 'bookmarks', this.json('bookmarks'));
		},
		'editor.saved': function(key, text) {
			this.set(key, JSON.parse(text));
		},
		/*
		 * TODO: Implement a control with handles and a cancel button so
		 * that you can fine-tune mouse selection.
		 */
		mouse: function(ractive_event) {
			var self;
			var canvas, ctx, event, x0, x1, y0, y1;

			self = this;
			event = ractive_event.original;

			if(event.button !== 0) return;

			x0 = Math.floor(event.pageX - this.canvas.offset().left);
			y0 = Math.floor(event.pageY - this.canvas.offset().top);
			x1 = x0;
			y1 = y0;

			this.canvas.on('mousemove mouseup mouseout', function handler(e) {
				var cancel, clear, ch, cw, dr, di, h, i, r, v, w;

				cw = self.canvas.width();
				ch = self.canvas.height();

				clear = function() {
					self.ctx.save();
					self.ctx.clearRect(0, 0, cw, ch);
					self.ctx.restore();
				};

				cancel = function() {
					self.canvas.off('mousemove mouseup mouseout', handler);
					clear();
				};

				if(e.type === 'mouseout') {
					cancel();
					return;
				}

				if(e.type === 'mouseup' && x0 !== x1 && y0 !== y1) {
					cancel();

					v = self.get('view');

					r = v.rmin;
					i = v.imin;
					w = v.rmax - r;
					h = v.imax - i;
					dr = w / cw;
					di = h / ch;

					v.rmin = r + dr * x0;
					v.imin = i + di * y0;
					v.rmax = r + dr * x1;
					v.imax = i + di * y1;
					self.update('view');

					return;
				}

				x1 = Math.floor(e.pageX - self.canvas.offset().left);
				y1 = y0 + ((x1 - x0) * (ch / cw));

				clear();
				self.ctx.strokeStyle = "rgba(0, 0, 0, 0.80)";
				self.ctx.strokeRect(x0 - 1.5, y0 - 1.5, x1 - x0 + 1, y1 - y0 + 1);
				self.ctx.strokeStyle = "rgba(255, 255, 255, 0.80)";
				self.ctx.strokeRect(x0 - 0.5, y0 - 0.5, x1 - x0, y1 - y0);
				self.ctx.fillStyle = "rgba(0, 220, 255, 0.20)";
				self.ctx.fillRect(x0, y0, x1 - x0 - 1.5, y1 - y0 - 1.5);
			});
		}
	},
	json: function(key) {
		return JSON.stringify(this.get(key), null, 2);
	},
	view_url: function(name) {
		//return "/png?" + $.param(this.get('view')) + '&render_id=' + this.get('render_id');
		return "/png?" +
			"i=" + this.get("i") +
			"&e=" + this.get("e") +
			"&m=" + this.get("m") +
			"&c=" + this.get("c") +
			"&r=" + this.get("r") +
			"&s=" + this.get("s") +
			"&p=" + this.get("p") +
			"&rmin=" + this.get("rmin") +
			"&rmax=" + this.get("rmax") +
			"&imin=" + this.get("imin") +
			"&imax=" + this.get("imax") +
			"&render_id=" + this.get('render_id');
	},
	update_view: function() {
		var i, image, self;

		self = this;
		image = this.find('#image');

		this.set('view.w', parseInt(image.width));
		this.set('view.h', parseInt(image.height));

		i = new Image();
		i.onload = function() {
			image.css({
				'background-image': "url(" + self.view_url() + ")"
			});
			self.ring.hidden = true;
		};
		this.ring.hidden = false;
		i.src = this.view_url();
	},
	translate_view: function(r, i) {
		var view;

		view = this.get('view');

		view.rmin += r;
		view.imin += i;
		view.rmax += r;
		view.imax += i;

		this.update('view');
	},
	scale_view: function(factor) {
		var rw, iw, rmid, imid, view;

		view = this.get('view');

		rw = (view.rmax - view.rmin) / 2.0;
		iw = (view.imax - view.imin) / 2.0;
		rmid = view.rmin + rw;
		imid = view.imin + iw;

		view.rmin = rmid - (rw * factor);
		view.imin = imid - (iw * factor);
		view.rmax = rmid + (rw * factor);
		view.imax = imid + (iw * factor);

		this.update('view');
	},
	copy_view: function(src_key, dst_key) {
		this.set(dst_key, JSON.parse(this.json(src_key)));
	},
});

