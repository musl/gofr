import Ractive from "ractive";

import Gofr from "./Gofr";
import Editor from "./Editor";

import template from "./Browser.html";

export default Ractive.extend({
  template: template,
  data: function() {
    return {
      key: '',
      text: '',
      mode: 'text',
      theme: 'github',
    };
  },
  on: {
    complete: function() {
      this.dialog = this.find('dialog');
      this.editor = ace.edit(this.find('div.editor'));
      this.session = this.editor.getSession();
      this.doc = this.session.getDocument();

      this.editor.setTheme("ace/theme/" + this.get('theme'));
      this.session.setMode("ace/mode/" + this.get('mode'));
      this.editor.$blockScrolling = Infinity;
    },
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
    },
  }
});
