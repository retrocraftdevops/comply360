# -*- coding: utf-8 -*-
from odoo import models, fields, api


class CrmLead(models.Model):
    """Extend CRM Lead with Comply360 integration fields"""
    _inherit = 'crm.lead'

    # Comply360 integration fields
    x_comply360_registration_id = fields.Char(
        string='Comply360 Registration ID',
        help='UUID of the registration in Comply360 system',
        index=True,
        copy=False,
    )

    x_comply360_registration_type = fields.Selection([
        ('pty_ltd', 'Pty Ltd Company'),
        ('cc', 'Close Corporation'),
        ('business_name', 'Business Name'),
        ('vat', 'VAT Registration'),
    ], string='Registration Type', help='Type of business registration')

    x_comply360_registration_number = fields.Char(
        string='Registration Number',
        help='Official registration number from CIPC/SARS',
    )

    x_comply360_status = fields.Selection([
        ('draft', 'Draft'),
        ('submitted', 'Submitted'),
        ('under_review', 'Under Review'),
        ('approved', 'Approved'),
        ('rejected', 'Rejected'),
        ('completed', 'Completed'),
    ], string='Comply360 Status', help='Current status in Comply360')

    x_comply360_sync_date = fields.Datetime(
        string='Last Sync Date',
        help='Last synchronization date with Comply360',
        readonly=True,
    )

    x_comply360_document_count = fields.Integer(
        string='Document Count',
        help='Number of documents uploaded in Comply360',
    )

    x_comply360_verified = fields.Boolean(
        string='Documents Verified',
        help='Whether all documents have been verified in Comply360',
        default=False,
    )

    @api.model
    def create(self, vals):
        """Override create to log Comply360 integrations"""
        lead = super(CrmLead, self).create(vals)
        if lead.x_comply360_registration_id:
            lead.message_post(
                body=f"Lead created from Comply360 registration: {lead.x_comply360_registration_id}"
            )
        return lead

    def write(self, vals):
        """Override write to track Comply360 status changes"""
        if 'x_comply360_status' in vals:
            for lead in self:
                old_status = lead.x_comply360_status
                new_status = vals['x_comply360_status']
                if old_status != new_status:
                    lead.message_post(
                        body=f"Comply360 status changed: {old_status} → {new_status}"
                    )
        return super(CrmLead, self).write(vals)

    def action_view_comply360_portal(self):
        """Open Comply360 portal for this registration"""
        self.ensure_one()
        if not self.x_comply360_registration_id:
            return

        base_url = self.env['ir.config_parameter'].sudo().get_param('comply360.portal_url', 'http://localhost:3000')
        url = f"{base_url}/registrations/{self.x_comply360_registration_id}"

        return {
            'type': 'ir.actions.act_url',
            'url': url,
            'target': 'new',
        }
