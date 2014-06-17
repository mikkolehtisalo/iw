Ext.define('IW.controller.Wikis', {
    extend: 'Ext.app.Controller',
    views: [
    'wiki.List',
    'wiki.Edit',
    'wiki.Window'
    ],
    stores: [
    'Wikis',
    'UserGroupSearch',
    'FavoriteWikis'
    ],
    models: [
    'Wiki',
    'UserGroupSearchItem'
    ],
    init: function() {
        this.control({
            'wikiwindow': {
                itemnewbuttonclick: this.newWiki,
                itemdeleteclick: this.deleteWiki
            },
            'wikilist': {
                itemdblclick: this.openWiki,
                itemeditbuttonclick : this.editWiki,
                itemunfavoriteclick : this.unfavoriteWiki,
                itemfavoriteclick : this.favoriteWiki
            },
            'wikiedit button[action=save]': {
                click: this.updateWiki
            },
            'wikiedit': {
                searchSelect: this.searchSelect
            }
        });
    },
    createUserButton: function(name, hiddenInput) {
        var butt = Ext.create('Ext.button.Split', {
            text: name,
            cls: 'iwsplitbutton',
            margin: '0 5px 5px 0',
            handler: function() {},
            arrowHandler: function(button, event) {
                var aclString = hiddenInput.value.split(',');
                var idx = aclString.indexOf(button.text); 
                if(idx!=-1) aclString.splice(idx, 1); 
                button.destroy();
                hiddenInput.setValue(aclString.join());
            }
        });
        return butt;
    },
    searchSelect: function(combo, record, hidden, panel) {
        hiddenInput = Ext.getCmp(hidden);
        var aclString = hiddenInput.value;
        if (!aclString) {
            aclString = '';
        }
        var idString = record.data.Id;

        // Update the value to the hidden input
        if (aclString.length > 0 ) {
            aclString = aclString + ',' + idString;
        } else {
            aclString = idString;
        }
        
        hiddenInput.setValue(aclString);

        // Add new button
        var panelCmpt = Ext.getCmp(panel);
        var butt = this.createUserButton(idString, hiddenInput);
        panelCmpt.add(butt);
    },
    buildACLEditor: function(hiddenInput, panel) {
        if (hiddenInput && hiddenInput.value) {
            var aclString = hiddenInput.value.split(',');
            for (i in aclString) {
                var st = aclString[i];
                var butt = this.createUserButton(st, hiddenInput);
                panel.add(butt);
            }
        }
    },
    editWiki: function(grid, record) {
        var view = Ext.widget('wikiedit');
        view.down('form').loadRecord(record);

        var rhidden = Ext.getCmp('wiki-hidden-read');
        var rpanel = Ext.getCmp('wiki-acl-read');
        this.buildACLEditor(rhidden, rpanel);

        var whidden = Ext.getCmp('wiki-hidden-write');
        var wpanel = Ext.getCmp('wiki-acl-write');
        this.buildACLEditor(whidden, wpanel);

        var ahidden = Ext.getCmp('wiki-hidden-admin');
        var apanel = Ext.getCmp('wiki-acl-admin');
        this.buildACLEditor(ahidden, apanel);
    },
    newWiki: function(grid, record) {
        var view = Ext.widget('wikiedit');
    },
    updateWiki: function(button) {
        var win    = button.up('window'),
        form   = win.down('form'),
        record = form.getRecord(),
        values = form.getValues();
        if (record == null) {
            // New Wiki!
            var wiki = Ext.create('IW.model.Wiki', {
                Title: values.Title,
                Description: values.Description,
                Readacl: values.Readacl,
                Writeacl: values.Writeacl,
                Adminacl: values.Adminacl,
                MatchedPermissions: ['admin'],
                Favorite: false
            });
            this.getWikisStore().add(wiki);
        } else {
            record.set(values);
            
        }
        win.close();
        this.getWikisStore().sync();
    },
    unfavoriteWiki: function(grid, record) {
        if (record.get('Favorite') == true) {
            record.set('Favorite', false);
            st = this.getFavoriteWikisStore();
            st.remove(st.getById(record.id));
            st.sync();
        }
    },
    favoriteWiki: function(grid, record) {
        if (record.get('Favorite') == false) {
            record.set('Favorite', true);

            // Only Wiki_id matters here
            var newfav = Ext.create('IW.model.FavoriteWiki', {
                Wiki_id : record.id,
                Username : 'auto',
                Modified : 'now',
                Status : 'ACTIVE'
            });
            newfav.save();
        }
    },
    deleteWiki: function(grid, record) {
        var me = this;
        Ext.Msg.confirm('Delete wiki?', 'You are about the delete the wiki <strong>'+record.data.Title+'</strong>. Are you sure you want to do this?', function(button) {
            if (button === 'yes') {
                // Setup ID
                var deletedId = record.data.Wiki_id;

                // Remove from store
                me.getWikisStore().remove(record);
                me.getWikisStore().sync();

                // Close wiki windows
                me.closeAllWikiWindows(deletedId);
            } 
        });
    },
    closeAllWikiWindows: function(closeId) {
        // Closes wiki window and the associated pages
        Ext.WindowManager.each(function (item) {   
            if (item.record && item.record.data.Wiki_id) {
                if (item.record.data.Wiki_id==closeId) {
                    item.close();
                }
            }
        });
    },
    openWiki: function(grid,record) {
        // Open a new window for this wiki UNLESS one is already open!
        
        var currentId = record.data.Wiki_id;

        var alreadyOpen = false;
        Ext.WindowManager.each(function (item) {   
            if (item.record && item.record.data.Wiki_id) {
                if (item.record.data.Wiki_id==currentId) {
                    alreadyOpen = true;
                }
            }
        });

        if (! alreadyOpen) {
            var tree = Ext.widget('treewindow', {
                record: record
            });
        }
    }
});


