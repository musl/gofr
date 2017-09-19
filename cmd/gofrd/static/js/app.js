/* vim: set ft=javascript sw=2 ts=2 : */

var Gofr = {};

Gofr.storage = localStorage;

Gofr.helpers = Ractive.defaults.data;
Gofr.helpers.complex = function(r, i) {
  return r + (i < 0 ? "" : "+") + i + "i";
};

Gofr.uuid = function() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

Gofr.ModalEditor = Ractive.extend({
  template: '#editor-tmpl',
  data: function() {
    return { 
      key: '',
      text: '',
      mode: 'text',
      theme: 'github',
    };
  },
  oncomplete: function() {
    this.dialog = this.find('dialog');
    this.editor = ace.edit(this.find('div.editor'));
    this.session = this.editor.getSession();
    this.doc = this.session.getDocument();

    this.editor.setTheme("ace/theme/" + this.get('theme'));
    this.session.setMode("ace/mode/" + this.get('mode'));
    this.editor.$blockScrolling = Infinity;

    this.on({
      save: function() {
        this.dialog.close();
        this.set('text', this.doc.getValue());
        this.fire('saved', this.get('key'), this.get('text'));
      },
      revert: function() {
        this.doc.setValue(this.get('text'));
      },
      dismiss: function() {
        this.dialog.close();
      },
      edit: function(key, text) {
        this.set('key', key);
        this.set('text', text);
        this.doc.setValue(text);
        this.dialog.showModal();
      }
    });
  }
});
Ractive.components.editor = Gofr.ModalEditor;

Gofr.FractalBrowser = Ractive.extend({
  template: '#browser-tmpl',
  data: function() {
    return {
      width: 960,
      height: 960,
      view: {},
      view_history: [],
      view_history_max: 3,
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
  components: {
    editor: Gofr.ModalEditor
  },
  onrender: function() {
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
      $.each(marks, function(key, value) {
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

    this.canvas = $(this.find('canvas'));
    this.ctx = this.canvas[0].getContext('2d');

    this.ring = $(this.find('div#ring'));
    window.ring = this.ring;

    this.update_view();
  },
  oncomplete: function() {
    this.observe('view', function() {
      Gofr.storage.setItem('gofr.browser.view', this.json('view'));
    });

    this.observe('view_history', function() {
      Gofr.storage.setItem('gofr.browser.view_history', this.json('view_history'));
    });

    this.observe('bookmarks', function() {
      Gofr.storage.setItem('gofr.browser.marks', this.json('bookmarks'));
    });

    this.observe('render_id', function() {
      Gofr.storage.setItem('gofr.browser.render_id', this.json('render_id'));
    });

    this.on({
      move_up: function() {
        this.translate_view(0.0, -0.0625 * (this.get('view.imax') - this.get('view.imin')));
        this.update_view();
      },
      move_down: function() {
        this.translate_view(0.0, 0.0625 * (this.get('view.imax') - this.get('view.imin')));
        this.update_view();
      },
      move_left: function() {
        this.translate_view(-0.0625 * (this.get('view.rmax') - this.get('view.rmin')), 0.0);
        this.update_view();
      },
      move_right: function() {
        this.translate_view(0.0625 * (this.get('view.rmax') - this.get('view.rmin')), 0.0);
        this.update_view();
      },
      zoom_in: function() {
        this.scale_view(0.9);
        this.update_view();
      },
      zoom_in_4x: function() {
        this.scale_view(0.6);
        this.update_view();
      },
      zoom_out: function() {
        this.scale_view(1.1);
        this.update_view();
      },
      zoom_out_4x: function() {
        this.scale_view(1.4);
        this.update_view();
      },	
      update_view: function() {
        this.update_view();
      },
      go_back: function() {
        this.set('view', this.get('view_history').shift());
        this.update_view(true);
      },
      go_to_bookmark: function(event) {
        var bookmark, name;

        name = 'bookmarks.' + $(event.node).data('bookmark');
        if(!name in this.get('bookmarks')) { return; }
        this.copy_view(name, 'view');
        this.update_view();
      },
      add_bookmark: function(event) {
        var name;

        // TODO FIX THIS window.prompt 
        name = window.prompt('Enter a name:');

        if(name) {
          this.copy_view('view', 'bookmarks.' + name);
          this.set('bookmarks.' + name + '.editable', true); 
        }
      },
      update_bookmark: function(event) {
        var name;

        name = $(event.node).data('bookmark');
        this.copy_view('view', 'bookmarks.' + name);
        this.set('bookmarks.' + name + '.editable', true); 
      },
      delete_bookmark: function(event) {
        var bookmarks;

        bookmarks = this.get('bookmarks');
        delete bookmarks[$(event.node).data('bookmark')];
        this.update('bookmarks');
      },
      edit_view: function() {
        var editor;

        editor = this.findComponent('editor');
        editor.fire('edit', 'view', this.json('view'));
      },
      edit_bookmarks: function() {
        var editor;

        editor = this.findComponent('editor');
        editor.fire('edit', 'bookmarks', this.json('bookmarks'));
      },
      'editor.saved': function(key, text) {
        this.set(key, JSON.parse(text));
        this.update_view();
      },
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

            self.update_view();
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
    });
  },
  json: function(key) {
    return JSON.stringify(this.get(key), null, 2);
  },
  view_url: function(name) {
    return "/png?" + $.param(this.get('view')) + '&render_id=' + this.get('render_id');
  },
  update_view: function(dont_update_history) {
    var image, self;

    self = this;
    image = $('#image');

    // The view is always a square who's width is determined by the
    // parent element. TODO hook up resizing?
    this.set('view.w', parseInt(image.width()));
    this.set('view.h', parseInt(image.width()));

    i = new Image();
    i.onload = function() {
      image.css({
        'background-image': "url(" + self.view_url() + ")"
      });
      self.ring.hide();
    };
    this.ring.show();
    i.src = this.view_url();

    if(dont_update_history) return;
    this.update_history();
  },
  update_history: function () {
    vh = this.get('view_history');

    vh.unshift(JSON.parse(this.json('view')));
    while(vh.length > this.get('view_history_max')) {
      vh.pop();
    }
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
Ractive.components.browser = Gofr.FractalBrowser;

$(document).ready(function() {
  var ractive = Ractive({
    el: '#gofr',
    template: '#gofr-tmpl',
    data: function() { return {}; },
    components: {
      browser: Gofr.FractalBrowser,
    },
  });
});

