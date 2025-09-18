-- CREATE TRIGGER IF NOT EXISTS addMemberLoanInvoiceDetail AFTER INSERT ON memberLoanRepayment FOR EACH ROW BEGIN CREATE TEMP TABLE IF NOT EXISTS varsLoanRepayment (amount REAL);

-- INSERT INTO
--   varsLoanRepayment (amount)
-- VALUES
--   (NEW.amount);

-- INSERT INTO
--   memberLoanInvoiceDetail (memberLoanInvoiceId, loanComponent, billedAmount)
-- VALUES
--   (NEW.id, "Interest", NEW.interest),
--   (NEW.id, "Instalment", NEW.instalment),
--   (NEW.id, "Insurance", NEW.insurance);

-- INSERT INTO
--   memberLoanInvoiceDetail (memberLoanInvoiceId, loanComponent, billedAmount)
-- SELECT
--   NEW.id,
--   "Processing Fee",
--   NEW.processingFee
-- WHERE
--   NEW.processingFee > 0;

-- DROP TABLE IF EXISTS varsLoanRepayment;

-- END;