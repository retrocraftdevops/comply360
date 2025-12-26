# -*- coding: utf-8 -*-
{
    'name': 'Comply360 Branding',
    'version': '19.0.1.0.0',
    'category': 'Hidden',
    'summary': 'Remove Odoo branding and apply Comply360 branding',
    'description': """
        Comply360 Custom Branding
        =========================
        This module removes Odoo branding and applies Comply360 branding throughout the interface.

        Features:
        - Removes "Powered by Odoo" footer
        - Customizes login page
        - Updates page titles
        - Applies Comply360 branding
    """,
    'author': 'Comply360',
    'website': 'https://comply360.com',
    'license': 'LGPL-3',
    'depends': ['base', 'web'],
    'data': [
        'views/webclient_templates.xml',
    ],
    'assets': {
        'web.assets_frontend': [
            'comply360_branding/static/src/css/branding.css',
        ],
        'web.assets_backend': [
            'comply360_branding/static/src/css/branding.css',
        ],
    },
    'installable': True,
    'auto_install': True,
    'application': False,
}
