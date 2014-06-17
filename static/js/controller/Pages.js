Ext.define('IW.controller.Pages', {
    extend: 'Ext.app.Controller',
    models: ['Page', 'ContentField', 'Attachment', 'Lock'],
    stores: ['Pages', 'ContentFields', 'Attachments', 'Locks'],
    views: ['page.PageTree','page.PageWindow','page.TreeWindow','page.Access'],
    requires: ['IW.Utilities'],

    init: function() {
        this.control({
            'pagetree': {
                itemdblclick: this.openPage,
                select: this.selectPage,
                deselect: this.deselectPage
            },
            'pagewindow': {
                itemeditbuttonclick: this.editPage,
                itemsavebuttonclick: this.savePage,
                resize: this.resizeWindow,
                itemrefreshbuttonclick: this.refreshPage,
                itemcancelbuttonclick: this.cancelEdit,
                close: this.closeWindow,
                pagelinkclick: this.openPageLink
            },
            'treewindow': {
                addpagebuttonclick: this.addPage,
                deletepagebuttonclick: this.deletePageConfirm,
                refreshtree: this.refreshTree,
                pageaccessrightsbuttonclick: this.showacldialog
            },
            'pageaccess': {
                searchSelect: this.searchSelect,
                saveACL: this.saveACL
            }
        });
    },
    saveACL: function (win) {
        var rhidden = Ext.getCmp('page-hidden-read');
        win.record.set('Readacl', rhidden.value);
        var whidden = Ext.getCmp('page-hidden-write');
        win.record.set('Writeacl', whidden.value);
        var ahidden = Ext.getCmp('page-hidden-admin');
        win.record.set('Adminacl', ahidden.value);

        var inherit = Ext.getCmp('stop-acl-inheritation');
        win.record.set('Stopinheritation', inherit.value);

        for (index in win.record.stores) {
            var store = win.record.stores[index];
            store.sync();
        }
      
        win.destroy();
    },
    searchSelect: function(combo, record, hidden, panel) {
        this.getController('Wikis').searchSelect(combo, record, hidden, panel);
    },
    showacldialog: function(event, target, owner, tool) {
        var rec = owner.up('window').down('panel').getSelectionModel().getSelection()[0];

        var dialog = Ext.create('IW.view.page.Access', {
            title: 'Edit access for '+rec.data.Title,
            record: rec
        });
        dialog.show();

        var rhidden = Ext.getCmp('page-hidden-read');
        rhidden.value = rec.data.Readacl;
        var rpanel = Ext.getCmp('page-acl-read');
        this.getController('Wikis').buildACLEditor(rhidden, rpanel);

        var whidden = Ext.getCmp('page-hidden-write');
        whidden.value = rec.data.Writeacl;
        var wpanel = Ext.getCmp('page-acl-write');
        this.getController('Wikis').buildACLEditor(whidden, wpanel);

        var ahidden = Ext.getCmp('page-hidden-admin');
        ahidden.value = rec.data.Adminacl;
        var apanel = Ext.getCmp('page-acl-admin');
        this.getController('Wikis').buildACLEditor(ahidden, apanel);

        var inherit = Ext.getCmp('stop-acl-inheritation');
        inherit.setValue(rec.data.Stopinheritation);
    },
    getAllChildren: function (node) { 
            var me = this;
            var all = new Array(); 

            if(!Ext.value(node,false)) { 
                    return []; 
            } 
            
            if(!node.hasChildNodes()) { 
                    return node;
            } else { 
                    all.push(node); 
                    node.eachChild(function (Mynode) { all = all.concat(me.getAllChildren(Mynode)); });         
            } 
            return all; 
    },
    refreshTree: function(event,target,owner,tool) {
        var tree = owner.up('window').down('treepanel');
        tree.getStore().reload(
                {
                    callback: function(records, options, success) {
                        // Nothing, see savePage()
                    }
                }
            );
    },
    deletePage: function (record, store) {
        // This is slightly bad, but TreeStore doesn't have method for deleting records.
        // It has to be done via model, which has to get the Proxy first from somewhere.
        // ....
        // Note that autoSync is not required, sync will occur in any case.
        record.setProxy(store.getProxy());
        var children = this.getAllChildren(record);
        Ext.Array.each(children, function(item, index, allchildren) {
            item.setProxy(store.getProxy());
        });

        record.destroy(false);
    },
    deletePageConfirm: function (event, target, owner, tool) {
        var rec = owner.up('window').down('panel').getSelectionModel().getSelection()[0];
        var st = owner.up('window').down('panel').store;
        var me = this;
        Ext.Msg.confirm('Delete page?', 'You are about the delete the page <strong>'+rec.data.Title+'</strong>. Are you sure you want to do this?', function(button) {
            if (button === 'yes') {
                me.deletePage(rec, st);
                me.closeAllChildren(rec.data.Path);
            } 
        });
    },
    closeAllChildren: function(path) {
        // Closes page and its children - using path
        Ext.WindowManager.each(function (item) {   
            if (item.record && item.record.data.Page_id && item.record.data.Path) {
                if (item.record.data.Path.indexOf(path) == 0) {
                    item.close();
                }
            }
        });
    },
    selectPage: function (panel, record, index, eOpts ) {
        var minus = panel.view.up('window').tools.Minus;
        var plus = panel.view.up('window').tools.Plus;
        var key = panel.view.up('window').tools.Key;

        if (minus.hidden && IW.Utilities.canwrite(record)) {
            minus.show();
        }
        if (key.hidden && IW.Utilities.canadmin(record)) {
            key.show();
        }
        if (IW.Utilities.canwrite(record)) {
            plus.show();
        } else {
            plus.hide();
        }
    },
    deselectPage: function (panel, record, index, eOpts ) {
        var minus = panel.view.up('window').tools.Minus;
        var plus = panel.view.up('window').tools.Plus;
        var key = panel.view.up('window').tools.Key;
        var wiki = panel.view.up('window').record;

        if (!minus.hidden) {
            minus.hide();
        }
        if (!key.hidden) {
            key.hide();
        }
        if (IW.Utilities.canwrite(wiki)) {
            plus.show();
        } else {
            plus.hide();
        }

    },
    openPageLink: function(win, element) {
        // Fake something that will work with openPage, reusing the same store all over again
        //var store = win.store;

        var wiki_id = element.id.split('/')[0]
        var parent_id = element.id.split('/')[1]
        var page_id = element.id.split('/')[2]

        var treeStore = Ext.create('IW.store.Pages');
        if (parent_id != "") {
            treeStore.setRootNode({
                expanded: true,
                Wiki_id: wiki_id,
                Page_id: parent_id
            });
        } else {
            treeStore.setRootNode({
                expanded: true,
                Wiki_id: wiki_id
            });
        }

        treeStore.load({
            scope: this,
            callback: function(records, operation, success) {
                var pagerecord;

                var root = treeStore.getRootNode();

                // treeStore doesn't have proper search method, yet
                // cascadeBy walks the tree until we get match
                root.cascadeBy(function (childNode) {
                    if (childNode.get('Page_id') == page_id)
                    {
                        pagerecord = childNode;
                    }
                }, this);

                // If we found one, open it!
                if (pagerecord) {
                    var panel = {
                        store: treeStore
                    };

                    this.getController('Pages').openPage(panel, pagerecord, null, null, null, null);
                }
            }
        });
    },
    addPageLinkEvents: function(win) {
        var ow = win;

        /*
        Ext.util.Observable.capture(ow, function(){
            console.log(arguments);
        });
        */
        ow.down('panel').body.on('click', function (evt, el, o) {
                ow.fireEvent('pagelinkclick', ow, el);
            }, this, { 
                delegate : '.extpagelink'
            });
    },
    fixPageAnchorLinks: function(win) {
        // First we will add the TOC links of window a random prefix, then the targets. This will keep
        // the anchors of pages separated. However, this will have to be noted when adding "open & edit"
        // links or similar to the subchapters - there will be extra nonsense on the ids.
        var randomCode = Math.random().toString(36).substring(7);

        var linkElems = win.down('panel').getEl().select("a.extpageanchor").elements;
        if (linkElems && linkElems.length > 0) {
            Ext.Array.forEach(linkElems, function (item, index, allItems) {
                var link = item.href.split('#');
                item.href = '#' + randomCode + '+' + link[1];
            });
            // Now change the titles from below...
            var titleElems = win.down('panel').getEl().select('.extpagetarget').elements;
            if (titleElems && titleElems.length > 0) {
                Ext.Array.forEach(titleElems, function (item, index, allItems) {
                    item.id = randomCode + '+' + item.id;
                });
            }
        }
    },
    openPage: function(panel, record, item, index, e, eOpts ) {
        console.log("openPage");
        var pw = Ext.widget('pagewindow', {
            title: record.data.Title,
            store: panel.store,
            record: record
        });
        pw.show();
        this.refreshPage(null, null, pw, null);
    },
    refreshPage: function(event, target, owner, tool) {
        var win;
        var me = this;
        if (owner.xtype == 'pagewindow' ) {
            win = owner;
        } else {
            win = owner.up('window');
        }
        var record = win.record;
        if (record) {
            // Reload and update
            var ContentField = Ext.ModelManager.getModel('IW.model.ContentField');
            var fieldId = record.data.Wiki_id.toString() + '/' + record.data.Page_id.toString(); 
            ContentField.load(fieldId, {
                success: function(field) {
                    win.down('panel').update(field.data.Contentwithmacros);
                    win.contentfield = field;
                    me.addPageLinkEvents (win);
                    me.fixPageAnchorLinks (win);
                }, 
                failure: function(field) {
                    console.log('Unable to load the content field');
                }
            });
        } else {
            console.log('Refresh requested but record was undefined!');
        }
    },
    cancelEdit: function(event, target, owner, tool, blockrefresh) {
        var panel = owner.up('window').down('panel');
        var editor = CKEDITOR.instances[panel.body.down('div').id];
        var record = owner.up('window').record;
        // Same...
        owner.up('window').tools.Refresh.show();
        owner.up('window').tools.Edit.show();
        owner.up('window').tools.Save.hide();
        owner.up('window').tools.Cancel.hide();
        editor.destroy();
        owner.up('window').setAutoScroll(true);
        panel.editorenabled = false;
        // Reset content field
        var me = this;
        setTimeout(function() {
            // If we reload the content faster than it has been saved, there will be problems...
            // This causes flicker, but the content shows more reliably. Should take an other look
            // later.
            me.refreshPage(null, null, owner.up('window'), null);
        }, 1000); 
        // Restore title
        owner.up('window').getHeader().remove(0);
        owner.up('window').setTitle(record.data.Title);
        // Remove listener that can prevent dragging
        owner.up('window').dd.removeListener('beforedragstart');
        // Remove the locks
        this.unlockPage(record.data.Page_id, record.data.Wiki_id);
    },
    editPage: function(event, target, owner, tool) {
        // Used to check for locks
        var me = this;
        var wind = owner.up('window');
        var page_id = wind.record.data.Page_id;
        var wiki_id = wind.record.data.Wiki_id;

        var Lock = Ext.ModelManager.getModel('IW.model.Lock');

        Lock.load(wiki_id+'/'+page_id, {
            success: function(lock) {
                console.log(lock);
                if (lock.data.Target_id != '') {
                    Ext.Msg.alert('Potential edit conflict!', 'User <em>'+ lock.data.Realname +
                        "</em> started editing this content at " + lock.data.Modified +
                        ". You may continue, but saving might cause conflicts.", function() {
                            // Open for editing
                            me.editPageProceed(event, target, owner, tool);
                        }, this );
                } else {
                    // There is probably no lock, so open the page for editing and lock
                    me.lockPage(page_id, wiki_id);
                    me.editPageProceed(event, target, owner, tool);
                }

            },
            failure: function(lock) {
                // There is probably no lock, so open the page for editing and lock
                me.lockPage(page_id, wiki_id);
                me.editPageProceed(event, target, owner, tool);
            }
        });
    },
    lockPage: function(page_id, wiki_id) {
        var lock = Ext.create('IW.model.Lock');
        lock.data.Wiki_id = wiki_id;
        lock.data.Target_id = page_id;
        lock.data.id = wiki_id+'/'+page_id;
        lock.save();
    },
    unlockPage: function(page_id, wiki_id) {
        var Lock = Ext.ModelManager.getModel('IW.model.Lock');
        Lock.load(wiki_id+'/'+page_id, {
            success: function (lock) {
                if (lock.get('Wiki_id') != "") {
                    // Destroy the lock
                    lock.destroy();
                }
            },
            failure: function(lock) {
                // Either not locked or error
            }
        });
    },
    editPageProceed: function(event, target, owner, tool) {
        var window = owner.up('window');
        // Change the content to version without macros
        window.down('panel').update(window.contentfield.data.Content);
        // window.down('panel').update(window.record.data.content);
        // Change the header to edit mode
        //-------------------------------------
        var killDrag = false;

        // window.dd is an Ext.util.ComponentDragger
        var dragEvent = window.dd.on({
            // beforedragstart event can cancel the drag
            beforedragstart: function(dd, e) {
                if (killDrag) {
                    return false;
                }
            }
        });

        var header = window.getHeader();
        var origTitle = header.title;
        header.setTitle(''); // Get rid of the original for now
        var field = Ext.create('Ext.form.field.Text', {
            name: 'Title',
            allowBlank: false,
            cls: 'pagetitleinput',
            value: origTitle,
            listeners : {
                el : {
                    delegate : 'input',
                    mouseout: function() {
                        killDrag = false;
                    },
                    mouseover: function() {
                        killDrag = true;
                    }
                }
            }
        });
        header.insert(0, field); // First, before the tools (several buttons there)

        // Change to content area to edit mode
        // -----------------------------------
        var panel = owner.up('window').down('panel');
        var contentdiv = panel.body.down('div');

        // Remove scrollbars, editor will handle this thing
        owner.up('window').setAutoScroll(false);

        // Boolean variable for other functions
        panel.editorenabled = true;
        // Enable the editor
        var newID = 'panel-'+owner.up('window').record.data.Page_id+'-innerCt'; // Generate ID that will be always the same for editing the same content
        contentdiv.dom.id = newID;
        contentdiv.id = newID;
        var editor = CKEDITOR.replace( contentdiv.id , { 
            baseFloatZIndex: 99999, 
            on: {
                // Resize once to fit the panel
                instanceReady: function(evt) {
                    var editor = evt.editor;
                    var width = panel.getWidth();
                    var height = panel.getHeight();
                    editor.resize(width,height);
                    editor.wiki = window.record.data.Wiki_id;
                }
            }
        });  

        // Change the tools
        owner.up('window').tools.Refresh.hide();
        owner.up('window').tools.Edit.hide();
        owner.up('window').tools.Save.show();
        owner.up('window').tools.Cancel.show();
    },
    savePage: function(event, target, owner, tool) {
        var panel = owner.up('window').down('panel');
        var editor = CKEDITOR.instances[panel.body.down('div').id];
        var record = owner.up('window').record;
        var contentfield = owner.up('window').contentfield;
        var newtitle = owner.up('window').getHeader().getComponent(0).value;

        // editor.checkDirty() failed for unknown reason, always reported "false" after setData()
        // Handle the content area
        if (contentfield.data.Content!=editor.getData()) {
            contentfield.data.Content=editor.getData();
            var fieldId = record.data.Wiki_id.toString() + '/' + record.data.Page_id.toString();

            contentfield.setId(fieldId);
            contentfield.save({ 
                success: function(record, operation) {
                    if (newtitle==record.data.Title) {
                        // Nothing so far
                    }
                },
                failure: function(record, operation) {
                    console.log('Unable to save content field!');
                }
            });
        }

        // Handle the title
        if (newtitle!=record.data.Title) {
            var store = owner.up('window').store;

            //var target = store.getById(record.data.internalId);
            //console.log(target);
            //record.data.title = newtitle;
            record.set('Title', newtitle);
            owner.up('window').record = record;

            // Save the title
            // The trick we did with separate stores will cause the following to be necessary.
            // There might be a cleaner solution to this though, have to investigate it later.
            var me = this;
            for (index in record.stores) {
                st = record.stores[index];
                // This is a bit messy, but we have to asynchronously sync (to save changes),
                // reload the item back from store (to get macros calculated), and set the content
                // to the reloaded...
                st.sync({
                    success: function() {
                        // Reload
                        var newstore = null;
                        // No idea why, but sometimes it's tree.treeStore, and sometimes treeStore
                        if (st.treeStore) {
                            newstore = st.treeStore;
                        } else {
                            newstore = st.tree.treeStore;
                        }

                        newstore.reload({
                            callback: function(records, options, success) {
                                //target = st.getById(record.data.id);
                                // Set the content item to record
                                //console.log(target);
                                //owner.up('window').record = target;
                                // Practically the same after everything has been saved successfully...
                                me.cancelEdit(event, target, owner, tool);
                            }
                        });
                    },
                    failure: function() {
                        console.log('Unable to save edited data!');
                    },
                    scope: this
                });
            } 
        } else {
            this.cancelEdit(event, target, owner, tool);
        }
    },
    closeWindow: function(window, eOpts) {
        // Remove locks
        var record = window.record;
        this.unlockPage(record.data.Page_id, record.data.Wiki_id);

        // Destroy the editors
        var panel = window.down('panel');
        if (panel.editorenabled) {
            var editor = CKEDITOR.instances[panel.body.down('div').id];
            editor.destroy();
        }
        window.down('panel').destroy();
        window.destroy();
    },
    resizeWindow: function ( window, width, height, oldWidth, oldHeight, eOpts ) {
        var panel = window.down('panel');
        if (panel.editorenabled) {
            var editor = CKEDITOR.instances[panel.body.down('div').id];
            editor.resize(panel.getWidth(),panel.getHeight());
        }

    },
    addPage: function (event, target, owner, tool) {
        var window = owner.up('window');
        var store = window.down('panel').getStore();
        var rootNode = store.getRootNode();

        // Determine the parent path. If nothing is selected, use rootNode.
        var parent = window.down('panel').getStore().getRootNode();
        var selected = window.down('panel').getSelectionModel().getSelection();
        if (selected.length > 0 ) {
            parent = selected[0];
        }
        var getPath = function(r,u,w) {
            if (r.data.Path && r.data.Path.length > 0) {
                return r.data.Path + '/' + u;
            } else {
                return u;
            }
        };

        // Create new model
        var uuidstr = Ext.data.IdGenerator.get('uuid').generate();
        var now = new Date();
        var newPage = Ext.create('IW.model.Page', {
            Page_id: uuidstr,
            Wiki_id: window.record.data.Wiki_id,
            Path: getPath(parent, uuidstr, window),
            Title: 'New Page',
            Create_user: '',
            Content: 'Insert content here',
            ContentWithMacros: 'Insert content here',
            Readacl: '',
            Writeacl: '',
            Adminacl: '',
            Modified: now.toJSON(),
            Stopinheritation: false,
            ChildNodes: [],
            loaded: true,
            MatchedPermissions: ['admin'],
            parentId: parent.internalId,
            root: rootNode,
        });

        // Add to ui
        parent.appendChild(newPage, false, true);
        parent.expand();
        window.down('panel').getSelectionModel().select(newPage);

        // Save the new page to store and open the window
        newPage.setProxy(store.getProxy());
        var me = this;
        newPage.save({
            success: function (record, operation) {
                // Open the page
                me.openPage (window.down('panel'), newPage);
            },
            failure: function (record, operation) {
                console.log('Unable to save new page!');
            }
        });
    }
});
