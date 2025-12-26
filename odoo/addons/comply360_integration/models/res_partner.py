# -*- coding: utf-8 -*-
from odoo import models, fields, api


class ResPartner(models.Model):
    """Extend Partner with Comply360 integration fields"""
    _inherit = 'res.partner'

    # Comply360 integration fields
    x_comply360_client_id = fields.Char(
        string='Comply360 Client ID',
        help='UUID of the client in Comply360 system',
        index=True,
        copy=False,
    )

    x_comply360_agent_id = fields.Char(
        string='Comply360 Agent ID',
        help='UUID of the agent in Comply360 system (for commission partners)',
        index=True,
        copy=False,
    )

    x_is_comply360_agent = fields.Boolean(
        string='Is Comply360 Agent',
        help='Indicates if this partner is a Comply360 agent receiving commissions',
        default=False,
    )

    x_comply360_registration_numbers = fields.Text(
        string='Registration Numbers',
        help='List of registration numbers managed through Comply360',
    )

    x_comply360_total_commissions = fields.Monetary(
        string='Total Commissions',
        help='Total commissions earned (for agents)',
        currency_field='currency_id',
        compute='_compute_comply360_commissions',
        store=True,
    )

    x_comply360_pending_commissions = fields.Monetary(
        string='Pending Commissions',
        help='Pending commission amount (for agents)',
        currency_field='currency_id',
        compute='_compute_comply360_commissions',
        store=True,
    )

    x_comply360_sync_date = fields.Datetime(
        string='Last Sync Date',
        help='Last synchronization date with Comply360',
        readonly=True,
    )

    x_comply360_commission_ids = fields.One2many(
        'x_commission',
        'partner_id',
        string='Commissions',
        help='Commission records from Comply360',
    )

    @api.depends('x_comply360_commission_ids', 'x_comply360_commission_ids.state', 'x_comply360_commission_ids.amount')
    def _compute_comply360_commissions(self):
        """Compute total and pending commissions"""
        for partner in self:
            total = sum(partner.x_comply360_commission_ids.mapped('amount'))
            pending = sum(
                partner.x_comply360_commission_ids.filtered(
                    lambda c: c.state in ['draft', 'approved']
                ).mapped('amount')
            )
            partner.x_comply360_total_commissions = total
            partner.x_comply360_pending_commissions = pending

    def action_view_comply360_commissions(self):
        """View all commissions for this agent"""
        self.ensure_one()
        return {
            'name': 'Commissions',
            'type': 'ir.actions.act_window',
            'res_model': 'x_commission',
            'view_mode': 'tree,form',
            'domain': [('partner_id', '=', self.id)],
            'context': {'default_partner_id': self.id},
        }
