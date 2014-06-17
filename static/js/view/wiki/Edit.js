Ext.define('IW.view.wiki.Edit', {
    extend: 'Ext.window.Window',
    alias: 'widget.wikiedit',

    title: 'Edit Wiki',
    layout: 'fit',
    width: 600,
    height: 400,
    autoShow: true,
    border: 0,
    modal: true,

    initComponent: function() {
        this.items = [
            {
                xtype: 'form',
                layout: 'anchor',
                overflowY: 'auto',
                padding: '5px 5px 5px 5px',
                border: 0,
                items: [
                    {
                        xtype: 'textfield',
                        name : 'Title',
                        fieldLabel: 'Title',
                        allowBlank: false,
                        minLength: 4,
                        maxLength: 100,
                        anchor: '100%'
                    },
                    {
                        xtype: 'textareafield',
                        name: 'Description',
                        fieldLabel: 'Description',
                        maxLength: 1024,
                        anchor: '100%'
                    },
                    {
                        xtype: 'hiddenfield',
                        id: 'wiki-hidden-read',
                        name : 'Readacl',
                        fieldLabel: 'Read',
                    },
                    {
                        xtype: 'hiddenfield',
                        id: 'wiki-hidden-write',
                        name : 'Writeacl',
                        fieldLabel: 'Write',
                    },
                    {
                        xtype: 'hiddenfield',
                        id: 'wiki-hidden-admin',
                        name : 'Adminacl',
                        fieldLabel: 'Admin',
                    },
                    {
                        xtype: 'panel',
                        layout: 'hbox',
                        anchor: '100%',
                        padding: '0 0 5px 0',
                        border: 0,
                        items: [{
                            xtype: 'panel',
                            html: 'Readers:',
                            width: 105,
                            border: 0
                        }, {
                            xtype: 'panel',
                            id: 'wiki-acl-read',
                            flex: 1,
                            border: 0
                        }]
                    },
                    {
                            xtype: 'combobox',
                            fieldLabel: 'Add reader',
                            store: 'UserGroupSearch',
                            displayField: 'Name',
                            valueField: 'Id',
                            anchor: '100%',
                            listeners: {
                                select: function(combo, records, eOpts) {
                                    if (records && records.length > 0) {
                                        var record = records[0];
                                        this.up('window').fireEvent('searchSelect', combo, record, 'wiki-hidden-read', 'wiki-acl-read');
                                    }
                                }
                            }
                    },
                    {
                        xtype: 'panel',
                        layout: 'hbox',
                        anchor: '100%',
                        padding: '0 0 5px 0',
                        border: 0,
                        items: [{
                            xtype: 'panel',
                            html: 'Writers:',
                            width: 105,
                            border: 0
                        }, {
                            xtype: 'panel',
                            id: 'wiki-acl-write',
                            flex: 1,
                            border: 0
                        }]
                    },
                    {
                            xtype: 'combobox',
                            fieldLabel: 'Add writer',
                            store: 'UserGroupSearch',
                            displayField: 'Name',
                            valueField: 'Id',
                            anchor: '100%',
                            listeners: {
                                select: function(combo, records, eOpts) {
                                    if (records && records.length > 0) {
                                        var record = records[0];
                                        this.up('window').fireEvent('searchSelect', combo, record, 'wiki-hidden-write', 'wiki-acl-write');
                                    }
                                }
                            }
                    },
                    {
                        xtype: 'panel',
                        layout: 'hbox',
                        anchor: '100%',
                        padding: '0 0 5px 0',
                        border: 0,
                        items: [{
                            xtype: 'panel',
                            html: 'Admins:',
                            width: 105,
                            border: 0
                        }, {
                            xtype: 'panel',
                            id: 'wiki-acl-admin',
                            flex: 1,
                            border: 0
                        }]
                    },
                    {
                            xtype: 'combobox',
                            fieldLabel: 'Add admin',
                            store: 'UserGroupSearch',
                            displayField: 'Name',
                            valueField: 'Id',
                            anchor: '100%',
                            listeners: {
                                select: function(combo, records, eOpts) {
                                    if (records && records.length > 0) {
                                        var record = records[0];
                                        this.up('window').fireEvent('searchSelect', combo, record, 'wiki-hidden-admin', 'wiki-acl-admin');
                                    }
                                }

                            }
                    }
                ]
            }
        ];

        this.buttons = [
            {
                text: 'Save',
                action: 'save'
            },
            {
                text: 'Cancel',
                scope: this,
                handler: this.close
            }
        ];

        this.callParent(arguments);
    }
});

