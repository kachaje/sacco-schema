INSERT INTO
  member (
    id,
    firstName,
    lastName,
    otherName,
    gender,
    title,
    maritalStatus,
    dateOfBirth,
    nationalIdentifier,
    utilityBillType,
    utilityBillNumber,
    phoneNumber
  )
VALUES
  (
    1,
    "Mary",
    "Banda",
    "",
    "Female",
    "Miss",
    "Single",
    "1999-09-01",
    "DHFYR8475",
    "ESCOM",
    "29383746",
    "09999999999"
  );

INSERT INTO
  memberContact (
    memberId,
    postalAddress,
    residentialAddress,
    homeVillage,
    homeTraditionalAuthority,
    homeDistrict
  )
VALUES
  (
    1,
    "P.O. Box 1000, Lilongwe",
    "Area 2, Lilongwe",
    "Songwe",
    "Kyungu",
    "Karonga"
  );
