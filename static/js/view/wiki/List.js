
Ext.define('IW.view.wiki.List' ,{
    extend: 'Ext.grid.Panel',
    alias: 'widget.wikilist',

    invalidateScrollerOnRefresh: false,
    viewConfig: {
        loadMask: false
    },
    hideHeaders: true,
    requires: [
        'Ext.grid.column.Action',
        'IW.Utilities'
    ],
    border: 0,

    store: 'Wikis',
    initComponent: function() {
        this.columns = [
            {header: 'Name', 
             dataIndex: 'Title',  
             flex: 1,
             renderer: function (value, metaData, record, row, col, store, gridView) {
               return '<p>'+String(value)+'</p><p><i>'+record.data.Description+'</i></p>';
             }
            },
            {
                xtype: 'actioncolumn',
                width: 60,
                items: [
                    {
                        iconCls: 'delete-col',
                        tooltip: 'Delete Wiki',
                        handler: function(grid, rowIndex, colIndex) {
                            var rec = grid.getStore().getAt(rowIndex);
                            this.up('window').fireEvent('itemdeleteclick', grid, rec);
                        },
                        isDisabled: function(view, rowIndex, colIndex, item, record) {
                            return ;
                        }
                    },
                    {
                    getClass: function(v, meta, rec) {
                        if (rec.get('Favorite')) {
                            return 'favorite-col';
                        } else {
                            return 'unfavorite-col';
                        }
                    },
                    getTip: function(v, meta, rec) {
                        if (rec.get('Favorite')) {
                            return 'Remove from favorites';
                        } else {
                            return 'Add to favorites';
                        }
                    },
                    handler: function(grid, rowIndex, colIndex) {
                        var rec = grid.getStore().getAt(rowIndex);
                        action = (rec.get('Favorite') ? 'Unfavorite' : 'Favorite');
                        if (action == 'Unfavorite') {
                            this.up('grid').fireEvent('itemunfavoriteclick', grid, rec);
                        } else {
                            this.up('grid').fireEvent('itemfavoriteclick', grid, rec);
                        }
                    }
                },
                {
                        iconCls: 'edit-col',
                        tooltip: 'Edit Properties',
                        handler: function(grid, rowIndex, colIndex) {
                            var rec = grid.getStore().getAt(rowIndex);
                            this.up('grid').fireEvent('itemeditbuttonclick', grid, rec);
                        },
                        isDisabled: function(view, rowIndex, colIndex, item, record) {
                            return !IW.Utilities.canadmin(record);
                        }
                }
                ]
            }
        ];

        this.callParent(arguments);
    }
});

