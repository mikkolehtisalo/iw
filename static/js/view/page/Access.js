Ext.define('IW.view.page.Access', {
    extend: 'Ext.window.Window',
    alias: 'widget.pageaccess',

    title: 'Edit access',
    layout: 'fit',
    width: 600,
    height: 350,
    autoShow: true,
    border: 0,
    modal: true,

    initComponent: function() {
        this.items = [
            {
                xtype: 'panel',
                layout: 'anchor',
                overflowY: 'auto',
                padding: '5px 5px 5px 5px',
                border: 0,
                items: [
                    {
                        xtype: 'hiddenfield',
                        id: 'page-hidden-read',
                        name : 'Readacl',
                        fieldLabel: 'Read',
                    },
                    {
                        xtype: 'hiddenfield',
                        id: 'page-hidden-write',
                        name : 'Writeacl',
                        fieldLabel: 'Write',
                    },
                    {
                        xtype: 'hiddenfield',
                        id: 'page-hidden-admin',
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
                            id: 'page-acl-read',
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
                                        this.up('window').fireEvent('searchSelect', combo, record, 'page-hidden-read', 'page-acl-read');
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
                            id: 'page-acl-write',
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
                                        this.up('window').fireEvent('searchSelect', combo, record, 'page-hidden-write', 'page-acl-write');
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
                            id: 'page-acl-admin',
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
                                        this.up('window').fireEvent('searchSelect', combo, record, 'page-hidden-admin', 'page-acl-admin');
                                    }
                                }

                            }
                    },
                    {
                        xtype: 'checkbox',
                        fieldLabel: 'Stop inheritation',
                        id: 'stop-acl-inheritation'
                    }
                ]
            }
        ];

        this.buttons = [
            {
                text: 'Save',
                handler: function() {
                    this.up('window').fireEvent('saveACL', this.up('window'));
                }
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

