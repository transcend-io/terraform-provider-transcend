resource "transcend_data_silo" "avc" {
  type                 = "promptAPerson"
  outer_type           = "coupa"
  notify_email_address = "dpo@coupa.com"
  description          = "Coupa is a cloud platform for business spend that offers a fully unified suite of financial applications for business spend management"
  is_live              = true
}