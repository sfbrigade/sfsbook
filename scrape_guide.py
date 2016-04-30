import bs4, json

def main():
   # open the guide html exported from Acrobat
   guide = open('RefGuide.html')
   # generate a Beautiful Soup object to operate on
   # I'm using lxml as the parser, 'pip install lxml' if you don't have it installed
   soup = bs4.BeautifulSoup(guide, 'lxml')

   # generate a list of all the headings for each entry in the guide, this is the most reliable pattern in the file
   # we'll use this to walk siblings to get the rest of the entry
   headings = soup.find_all('h1')
   # there should be 675 entries, some of these will just be numbers because junky formatting
   # we need to strip these out. There will still be some headings with numbers at the end, but that's ok
   # removing only number h1 entries allows us to use H1s as the dividing line for each section
   cleaned_headings = []
   for heading in headings:
      if heading.getText().isdigit() == False:
         cleaned_headings.append(heading)
   
   # Now that we have "cleaned" H1 entries, we can group all the still nosiy data between them
   # these are actual entry groupings, although they'll still need work after this
   # this should pull in all but 7 of the entries that have malformed lists in them, becuase of course they do
   # I'll add the missing 7 by hand and call it good enough. 
   entries = []
   for heading in cleaned_headings:
      this_entry = [heading]
      current = heading.find_next_sibling()
      while True:
         try:
            if current.name != 'h1':
               this_entry.append(current)
               current = current.find_next_sibling()
            elif current.name == 'h1' and current.getText().isdigit() == True:
               this_entry.append(current)
               current = current.find_next_sibling()
            else:
               entries.append(this_entry)
               break
         except AttributeError:
            print(" Needs to be entered by hand: " + str(heading))
            break
   # Strip out all the actual text from the rest of the cruft now that tags are no longer useful
   raw_guide_text = []
   for entry in entries:
      this_entry = []
      for line in entry:
         if line.getText().isprintable():
            this_entry.append(line.getText())
      raw_guide_text.append(this_entry)

   # Things we know for sure about the data in raw_guide_text at this point:
   # the entries are not all the same length, smallest one is 10 items, longest is 30
   # the first item in each entry will always be the title of the resource
   # the second item will be an address with just a few exceptions
   # the the third item is sometimes the address, but most often a number or the phone column label
   # phone numbers should always appear in the order "Crisis line: Business line: Fax: TDD:"
   # however, there isn't any indication which numbers are missing if there are < 4 in the set
   # the desctiption ALWAYS follows the 'Description:' entry
   # however, 14 of the entries don't have a 'Description:' field, but appear to have descriptions
   # some of the time the services offered follows the 'Services:' entry
   # however, 14 of the entries don't have a 'Services:' entry
   # and in 55 entries 'Wheelchair Accessible: Languages' followes the 'Services:' entry
   # Wheelchair Accessible, Languages, Populations served, and Categories seem to follow the same pattern as services and description
   # the website URL will always be in an entry that starts with 'Website'
   # the email address will always be in an entry that starts with 'Email'
   
   new_guide = []
   # We're popping the entries we can key with some certainty so we can
   # gather the rest of the stuff into a leftovers group
   for entry in raw_guide_text:
      new_entry = {}
      new_entry["name"] = entry.pop(0)
      # address may not be correct. This will just need human eyes on it
      new_entry["address"] = entry.pop(0)
      # if no description exists in the guide, we'll mark it "needs description"
      try:
         f = entry.index('Description:')
         new_entry["description"] = entry.pop(f+1)
         entry.pop(f)
      except ValueError:
         new_entry["description"] = "This entry needs a description"
      # if no services exists or is the weird junk data we'll mark it "needs services"
      # there is also a weird chance we might have the whole wheelchair entry in here,
      # if so, we'll just add it in right now and check for it in the wheelchar phase
      try:
         f = entry.index('Services:')
         if entry[f+1] == "Services: Wheelchair Accessible: Yes":
            new_entry["wheelchair"] = True
            new_entry["services"] = "This entry needs a services list"
            entry.pop(f+1)
            entry.pop(f)
         # these are the junk entries that might exist in this field
         elif entry[f+1] == "Wheelchair Accessible:": 
            new_entry["services"] = "This entry needs a services list"
            new_entry["wheelchair"] = "This entry needs Wheelchair accessibliy info"
            entry.pop(f+1)
            entry.pop(f)
         # this one tells us there isn't a languges list either
         elif entry[f+1] == "Wheelchair Accessible: Languages":
            new_entry["languges"] = "This entry needs a languages list"
            new_entry["services"] = "This entry needs a services list"
            new_entry["wheelchair"] = "This entry needs Wheelchair accessibliy info"
            entry.pop(f+1)
            entry.pop(f)
         else:
            new_entry["services"] = entry.pop(f+1)
            entry.pop(f)
      except ValueError:
          new_entry["services"] = "This entry needs a services list"
      # wheelchair accessibile it might already be filled in
      try:
         # if the wheelchair entry already exists this will be fine and we
         # can ignore it. Otherwise it will throw a KeyError
         new_entry["wheelchair"]
      except KeyError:
         try:
            f = entry.index('Wheelchair Accessible:')
            if entry[f+1] != "Yes":
               new_entry["wheelchair"] = "This entry needs Wheelchair accessibliy info"
               entry.pop(f+1)
               entry.pop(f)
            else:
               new_entry["wheelchair"] = entry[f+1]
               entry.pop(f)
         except ValueError:
            new_entry["wheelchair"] = "This entry needs Wheelchair accessibliy info"
      # languges it might already be filled in, just like wheelchair 
      try:
         new_entry["languges"]
      except KeyError:
         try:
            f = entry.index("Languages")
            if entry[f+1] == "Populations served:":
               new_entry["languges"] =  "This entry needs a languges list"
               new_entry["pops_served"] = "This entry needs a populations served list"
               entry.pop(f+1)
               entry.pop(f)
            elif entry[f+1] == "Website E-mail":
               new_entry["languges"] =  "This entry needs a languges list"
               new_entry["website"] = "This entry needs a website"
               new_entry["email"] = "This entry needs an email address"
               entry.pop(f+1)
               entry.pop(f)
            elif entry[f+1] == "Website":
               new_entry["languges"] =  "This entry needs a languges list"
               new_entry["website"] = "This entry needs a website"
               entry.pop(f+1)
               entry.pop(f)
            else:
               new_entry["languges"] = entry[f+1]
               entry.pop(f)
         except ValueError:
            new_entry["languges"] =  "This entry needs a languges list"
      # populations served, same checks as wheelchair
      try:
         new_entry["pops_served"]
      except KeyError:
         try:
            f = entry.index("Populations served:")
            if entry[f+1] == "Website E-mail":
               new_entry["pops_served"] =  "This entry needs a populations served list"
               new_entry["website"] = "This entry needs a website"
               new_entry["email"] = "This entry needs an email address"
               entry.pop(f+1)
               entry.pop(f)
            elif entry[f+1] == "Categories:":
               new_entry["pops_served"] =  "This entry needs a populations served list"
               new_entry["categories"] = "This entry needs a categories list"
               entry.pop(f+1)
               entry.pop(f)
            # there are a lot of web addresses hiding inside f+1, this pulls them out
            elif entry[f+1].find("Website") == 0:
               new_entry["website"] = entry.pop(f+1)[8:]
               new_entry["pops_served"] =  "This entry needs a populations served list"
               entry.pop(f)
            else:
               new_entry["pops_served"] = entry.pop(f+1)
               entry.pop(f)
         except ValueError:
            new_entry["pops_served"] =  "This entry needs a populations served list"
      # categories, same as populations et al.
      try:
         new_entry["categories"]
      except KeyError:
         try:
            f = entry.index("Categories:")
            # there are a handful of email addresses in f+1 across the set, some will
            # also include the categories after the email address, more cleanup by hand
            if entry[f+1].find("E-mail") == 0:
               new_entry["email"] = entry.pop(f+1)[7:]
               new_entry["categories"] = "This entry needs a categories list"
               entry.pop(f)
            else:
               new_entry["categories"] = entry.pop(f+1)
               entry.pop(f)
         except ValueError:
            new_entry["categories"] = "This entry needs a categories list"
      # website, same try except pattern as above
      try:
         new_entry["website"]
      except KeyError:
         # Website should only be in one place in the list, so if it showed up as
         # part of an early result, it should already be popped out of the list
         sub = "Website"
         search = [s for s in entry if sub in s]
         if search != []:
            new_entry["website"] = search[0][8:]
            # I don't need a try here because I know the search string exists
            entry.pop(entry.index(search[0]))
         else:
            new_entry["website"] = "This entry needs a website"
      # email, exact same as website, just looking for/trimming 'E-mail'
      try:
         new_entry["email"]
      except KeyError:
         # E-mail should only be in one place in the list, so if it showed up as
         # part of an early result, it should already be popped out of the list
         sub = "E-mail"
         search = [s for s in entry if sub in s]
         if search != []:
            new_entry["email"] = search[0][7:]
            # I don't need a try here because I know the search string exists
            entry.pop(entry.index(search[0]))
         else:
            new_entry["email"] = "This entry needs an email address"
      # everything else... I'm sorry future me and anyone else working on this
      # but this is all the stuff that is so inconsistant that I don't know what
      # else to do with it but lump it together and we'll have to have human eyes 
      # on it to figure out what it is. 
      new_entry["hand_sort"] = entry
      # add the new entry to new_guide
      new_guide.append(new_entry)

   # Now that new_guide is populated, convert it to JSON
   with open('refguide.json', 'w') as f:
      f.write(json.dumps(new_guide, sort_keys=True, indent=4))
   print ("All done!")

if __name__ == '__main__':
   main()