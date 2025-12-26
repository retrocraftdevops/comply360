# -*- coding: utf-8 -*-
from odoo import models, fields, api, _
from odoo.exceptions import ValidationError


class Commission(models.Model):
    """Commission tracking model for Comply360 integration"""
    _name = 'x_commission'
    _description = 'Comply360 Commission'
    _inherit = ['mail.thread', 'mail.activity.mixin']
    _order = 'date desc, id desc'

    name = fields.Char(
        string='Commission Reference',
        required=True,
        copy=False,
        default=lambda self: _('New'),
    )

    partner_id = fields.Many2one(
        'res.partner',
        string='Agent',
        required=True,
        domain=[('x_is_comply360_agent', '=', True)],
        tracking=True,
    )

    amount = fields.Monetary(
        string='Commission Amount',
        required=True,
        currency_field='currency_id',
        tracking=True,
    )

    currency_id = fields.Many2one(
        'res.currency',
        string='Currency',
        required=True,
        default=lambda self: self.env.company.currency_id,
    )

    date = fields.Date(
        string='Commission Date',
        required=True,
        default=fields.Date.context_today,
        tracking=True,
    )

    state = fields.Selection([
        ('draft', 'Draft'),
        ('approved', 'Approved'),
        ('paid', 'Paid'),
        ('refused', 'Refused'),
    ], string='Status', default='draft', required=True, tracking=True)

    description = fields.Text(
        string='Description',
        help='Commission description and details',
    )

    registration_reference = fields.Char(
        string='Registration Reference',
        help='Reference to the registration that generated this commission',
    )

    registration_type = fields.Selection([
        ('pty_ltd', 'Pty Ltd Company'),
        ('cc', 'Close Corporation'),
        ('business_name', 'Business Name'),
        ('vat', 'VAT Registration'),
    ], string='Registration Type')

    x_comply360_commission_id = fields.Char(
        string='Comply360 Commission ID',
        help='UUID of the commission in Comply360 system',
        index=True,
        copy=False,
    )

    x_comply360_registration_id = fields.Char(
        string='Comply360 Registration ID',
        help='UUID of the related registration in Comply360',
        index=True,
    )

    payment_date = fields.Date(
        string='Payment Date',
        readonly=True,
        tracking=True,
    )

    invoice_id = fields.Many2one(
        'account.move',
        string='Vendor Bill',
        help='Vendor bill for commission payment',
        readonly=True,
    )

    notes = fields.Html(string='Notes')

    @api.model_create_multi
    def create(self, vals_list):
        """Generate sequence on create"""
        for vals in vals_list:
            if vals.get('name', _('New')) == _('New'):
                vals['name'] = self.env['ir.sequence'].next_by_code('x_commission') or _('New')
        return super(Commission, self).create(vals_list)

    def action_approve(self):
        """Approve commission"""
        for commission in self:
            if commission.state != 'draft':
                raise ValidationError(_('Only draft commissions can be approved.'))
            commission.write({'state': 'approved'})
            commission.message_post(body=_('Commission approved'))

    def action_refuse(self):
        """Refuse commission"""
        for commission in self:
            if commission.state not in ['draft', 'approved']:
                raise ValidationError(_('Only draft or approved commissions can be refused.'))
            commission.write({'state': 'refused'})
            commission.message_post(body=_('Commission refused'))

    def action_mark_paid(self):
        """Mark commission as paid"""
        for commission in self:
            if commission.state != 'approved':
                raise ValidationError(_('Only approved commissions can be marked as paid.'))
            commission.write({
                'state': 'paid',
                'payment_date': fields.Date.context_today(self),
            })
            commission.message_post(body=_('Commission marked as paid'))

    def action_create_vendor_bill(self):
        """Create a vendor bill for this commission"""
        self.ensure_one()

        if self.invoice_id:
            raise ValidationError(_('A vendor bill already exists for this commission.'))

        if self.state != 'approved':
            raise ValidationError(_('Only approved commissions can generate vendor bills.'))

        # Create vendor bill
        invoice_vals = {
            'move_type': 'in_invoice',
            'partner_id': self.partner_id.id,
            'invoice_date': fields.Date.context_today(self),
            'ref': self.name,
            'invoice_line_ids': [(0, 0, {
                'name': f'Commission - {self.registration_reference or self.name}',
                'quantity': 1,
                'price_unit': self.amount,
            })],
        }

        invoice = self.env['account.move'].create(invoice_vals)
        self.invoice_id = invoice.id

        self.message_post(body=_('Vendor bill created: %s') % invoice.name)

        return {
            'type': 'ir.actions.act_window',
            'res_model': 'account.move',
            'res_id': invoice.id,
            'view_mode': 'form',
            'views': [(False, 'form')],
        }

    def action_view_invoice(self):
        """View related vendor bill"""
        self.ensure_one()
        if not self.invoice_id:
            raise ValidationError(_('No vendor bill exists for this commission.'))

        return {
            'type': 'ir.actions.act_window',
            'res_model': 'account.move',
            'res_id': self.invoice_id.id,
            'view_mode': 'form',
            'views': [(False, 'form')],
        }
