# -*- coding: utf-8 -*-
{
    'name': 'Comply360 Integration',
    'version': '1.0.0',
    'category': 'Integration',
    'summary': 'Integration module for Comply360 platform',
    'description': """
Comply360 Integration Module for Odoo 19
=========================================

This module extends Odoo 19 to integrate with the Comply360 compliance platform.

Features:
---------
* Custom fields on CRM leads for Comply360 registration tracking
* Custom commission tracking model
* Enhanced partner/customer management
* Webhook support for bi-directional sync
* Custom reports for compliance tracking

Dependencies:
------------
* crm (CRM module)
* sale (Sales module)
* account (Accounting module)
* contacts (Contacts module)
    """,
    'author': 'Comply360 Development Team',
    'website': 'https://comply360.com',
    'license': 'LGPL-3',
    'depends': [
        'base',
        'crm',
        'sale',
        'account',
        'contacts',
    ],
    'data': [
        'security/ir.model.access.csv',
        'views/crm_lead_views.xml',
        'views/res_partner_views.xml',
        'views/commission_views.xml',
        'views/menu_views.xml',
    ],
    'demo': [],
    'installable': True,
    'application': False,
    'auto_install': False,
}
