const Profile = () => {
  return (
    <div
      style={{
        marginTop: "6rem",
        textAlign: "center",
        display: "flex",

        gap: "3rem",
      }}
    >
      <User userName={`${user.username}`} isLoggedIn={true} name={user.name} />
      <Button className={classes.itemButton}>Follow</Button>
    </div>
  );
};