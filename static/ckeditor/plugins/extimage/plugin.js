Ext.define('CK.view.ImageList' ,{
    extend: 'Ext.grid.Panel',
    alias: 'widget.imagelist',

    hideHeaders: true,
    columns: [{ 
        text: 'Preview', 
        width: 84,
        dataIndex: 'Attachment_id',
        renderer: function(value) {
            return Ext.String.format('<img src="{0}" width="64" height="64"></img>', '/att/'+this.up('window').wiki+'/'+value+'?thumbnail=yes');
        }
    },{
        text: 'Filename',
        dataIndex: 'Filename',
        flex: 1
    },{
        text: 'Modified',
        dataIndex: 'Modified',
        width: 120
    },{
        xtype: 'actioncolumn',
        width: 40,
        items: [{
            iconCls: 'delete-col',
            tooltip: 'Delete Image',
            handler: function(grid, rowIndex, colIndex) {
                var store = grid.getStore();
                var rec = store.getAt(rowIndex);

                Ext.Msg.confirm('Delete the item?', 'Are you sure you want to delete the image '+rec.data.Filename, function(button) {
                    if (button === 'yes') {
                        // Will be wiki/attachment_id
                        rec.setId(rec.get('Attachment_id'));
                        store.remove(rec);
                        store.sync();
                    } else {
                        // Nothing really
                    }
                });
            }
        }]
    }],
    listeners: {
        itemdblclick: function(dv, record, item, index, e) {
            var zonepanel = Ext.getCmp('dropzonepanel');
            var win = zonepanel.up('window');
            var wiki = win.wiki;
            var editor = win.editor;
            // Add link to image
            var img = editor.document.createElement('img');
            img.setAttribute ('src','/att/'+wiki+'/'+record.data.Attachment_id);
            editor.insertElement (img);
            // Close window
            zonepanel.up('window').destroy();
        }
    }

});

Ext.define('CK.view.extimageWindow', {
    extend: 'Ext.window.Window',
    alias: 'widget.extimageWindow',

    height: 400,
    width: 500,
    title: 'Image selection',
    collapsible: false,
    modal: true,
    layout: 'fit',
    border: 0,
    defaults: {
        border: 0
    },
    items: [
    {
        xtype: 'panel',
        layout: 'border',
        items: [
        {
            region: 'center',
            xtype: 'tabpanel',
            layout: 'vbox',
            border: 0,
            items: [
            {
                title: 'Upload new image',
                margin: 10,
                defaults: {
                    border: 0
                },
                border: 0,
                items: [
                {
                    html: '<div id="drop_zone">Drop image file here</div>',
                    id: 'dropzonepanel',
                },{
                    xtype: 'panel',
                    id: 'orpanel',
                    html: 'or...',
                    padding: '10px 0 10px 0',
                },{
                    xtype: 'filefield',
                    id: 'filefieldpanel',
                    name: 'file',
                    width: 450,
                    margin: '5px 0 5px 0',
                    buttonText: 'Select Image...',
                    listeners: {
                        change: function(field, value, opts) { 
                            var zonepanel = Ext.getCmp('dropzonepanel');
                            var elem = document.getElementById(field.fileInputEl.id);
                            var file = elem.files[0];
                            zonepanel.up('window').down('panel').uploadNew(file, null);
                        }
                    }
                }
                ]
            },{
                title: 'Browse from gallery',
                border: 0,
                scroll: 'both',
                layout: 'fit',
                items: [{xtype: 'imagelist'}],
            },{
                title: 'Link external image',
                border: 0,
                layout: {
                    type: 'vbox',
                    align: 'center',
                    pack: 'center'
                },
                items: [{
                    width: 450,
                    height: 100,
                    border: 0,
                    layout: 'hbox',
                    items: [{
                        xtype: 'textfield',
                        id: 'imageUrlField',
                        width: 400,
                    },{
                        xtype: 'button',
                        margin: '0 0 0 5px',
                        text: 'Add',
                        handler : function() {
                            var url =  Ext.getCmp('imageUrlField').value;
                            var win = this.up('window');
                            var editor = win.editor;
                            // Add link to image
                            var img = editor.document.createElement('img');
                            img.setAttribute ('src', url);
                            editor.insertElement (img);
                            // Close window
                            win.destroy();
                        }
                    }
                    ]
                }]
            }
            ]
        }],
        listeners: {
            afterrender: function(win, opts) {
                var me = this;
                var dropElem = win.down('#dropzonepanel').getEl();
                dropElem.on('dragover', me.dragover);
                dropElem.on('drop', me.drop);
            }
        },
        dragover: function(evt, el, o) {
            evt.stopPropagation();
            evt.preventDefault();
        },
        uploadNew: function(file) {
            var zonepanel = Ext.getCmp('dropzonepanel');
            var win = zonepanel.up('window');
            var wiki = win.wiki;
            var editor = win.editor;
            var reader = new FileReader();

            reader.onloadend = function(evt) {
                var uuidstr = Ext.data.IdGenerator.get('uuid').generate();
                var now = new Date();
                var newAttachment = Ext.create('IW.model.Attachment', {
                    Attachment_id: uuidstr,
                    Wiki_id: wiki,
                    Attachment: evt.target.result,
                    Modified: now.toJSON(),
                    Filename: file.name,
                });
                newAttachment.proxy.url = '/api/attachments/'+wiki+'/'+uuidstr;
                newAttachment.save({ 
                    success: function(record, operation) {
                        // Add link to image
                        var img = editor.document.createElement('img');
                        img.setAttribute ('src','/att/'+wiki+'/'+uuidstr);
                        editor.insertElement (img);
                        // Close window
                        zonepanel.up('window').destroy();

                    },
                    failure: function(record, operation) {
                        //handle failure(s) here
                    }
                });
            };
            reader.readAsDataURL(file); // Safest with JSON
        },
        drop: function(evt, el, o) {
            evt.stopPropagation();
            evt.preventDefault();
            var cmp   = this.cmp,
            browserEvent = evt.browserEvent,
            dataTransfer = browserEvent.dataTransfer,
            files        = dataTransfer.files,
            numFiles     = files.length,
            file;
            var zonepanel = Ext.getCmp('dropzonepanel');

            // If dropped multiple, handle only the first file...
            if (numFiles > 0) {
                file = files[0];
                zonepanel.update ('<div id="drop_zone">'+file.name+'</div>');
                zonepanel.up('window').down('panel').uploadNew(file, evt);
            }
        }
    }
    ],
    initComponent: function() {
        this.callParent(arguments);
        // Can this be done with less code?
        var newAttachmentStore = Ext.create('IW.store.Attachments', {
                proxy: {
                type: 'rest',
                url: '/api/attachments/'+this.wiki+'/',
                reader: {
                    type: 'json',
                    successProperty: 'success'
                    }
                }
        });
        var grid = this.down('grid');
        grid.reconfigure(newAttachmentStore);
        grid.store.load();

        // Can this be done otherwise?
        /*
        grid.on('itemdblclick', function() { 
            console.log('lol')
        }, this);
*/
    }
});

Ext.define('CK.view.extimagePropertiesWindow', {
    extend: 'Ext.window.Window',
    alias: 'widget.extimagePropertiesWindow',

    height: 200,
    width: 400,
    title: 'Image properties',
    collapsible: false,
    modal: true,
    layout: 'fit',
    border: 0,
    defaults: {
        border: 0
    },
    items: [
        {
            xtype: 'panel',
            layout: 'border',
            padding: '5px 10px 5px 10px',
            items: [{
                region: 'center',
                xtype: 'panel',
                layout: 'vbox',
                border: 0,
                items: [{
                        xtype: 'textfield',
                        fieldLabel: 'Width',
                        id: 'propFormWidth',
                        width: 370,
                },{
                        xtype: 'textfield',
                        fieldLabel: 'Height',
                        id: 'propFormHeight', 
                        width: 370,
                },{
                        xtype: 'textfield',
                        fieldLabel: 'Description',
                        id: 'propFormDesc', 
                        width: 370,
                }
                ]
            }]
        }
    ],
    initComponent: function() {
        this.callParent(arguments);
        if (this.element.hasAttribute('width')) {
            Ext.getCmp('propFormWidth').setValue(this.element.getAttribute('width'));
        }
        if (this.element.hasAttribute('height')) {
            Ext.getCmp('propFormHeight').setValue(this.element.getAttribute('height'));
        }
        if (this.element.hasAttribute('alt')) {
            Ext.getCmp('propFormDesc').setValue(this.element.getAttribute('alt'));
        }
    },
    buttons: [
                {
                    text: 'Save',
                    handler: function() {
                        this.up('window').element.setAttribute('width', Ext.getCmp('propFormWidth').getValue());
                        this.up('window').element.setAttribute('height', Ext.getCmp('propFormHeight').getValue());
                        this.up('window').element.setAttribute('alt', Ext.getCmp('propFormDesc').getValue());
                        this.up('window').destroy();
                    }
                },
                {
                    text: 'Cancel',
                    handler: function() {
                        this.up('window').destroy();
                    }
                }
            ]
});


CKEDITOR.plugins.add( 'extimage', {
    icons: 'extimage',
    init: function( editor ) {

        //Plugin logic goes here.
        editor.addCommand( 'insertImage', {
            allowedContent: 'img[!src,alt,width,height]',
            exec: function( editor ) {
                var extWin = Ext.widget('extimageWindow', {
                    wiki: editor.wiki,
                    editor: editor
                });
                extWin.show();

                if (window.File && window.FileReader && window.FileList && window.Blob) {
                  //console.log('Great success! All the File APIs are supported.');
                } else {
                  //console.log('The File APIs are not fully supported in this browser.');
                }
            }
        });

        editor.ui.addButton( 'extimage', {
            label: 'Insert Image',
            command: 'insertImage',
            toolbar: 'insert'
        });

        // Context menu for image properties

        editor.addCommand( 'imageProperties', {
            exec: function (editor) {
                var selection = editor.getSelection();
                var element = selection.getSelectedElement();
                var propWin = Ext.widget('extimagePropertiesWindow', {
                    editor: editor,
                    selection: selection,
                    element: element
                });
                propWin.show();
            }
        });

        if ( editor.contextMenu ) {
            editor.addMenuGroup( 'imgGroup' );
            editor.addMenuItem( 'imgItem', {
                label: 'Image Properties',
                icon: this.path + 'icons/extimage.png',
                command: 'imageProperties',
                group: 'imgGroup'
            });

            editor.contextMenu.addListener( function( element ) {
                if ( element.getAscendant( 'img', true ) ) {
                    return { imgItem: CKEDITOR.TRISTATE_OFF };
                }
            });

        }

    }
});
